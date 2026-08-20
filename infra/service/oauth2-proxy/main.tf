terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 3.2"
    }
    vault = {
      source  = "hashicorp/vault"
      version = "~> 5.10"
    }
  }
}

module "secret" {
  source        = "./modules/secret"
  namespace     = var.namespace
  secret_name   = var.credentials_secret_name
  vault_kv_path = var.vault_kv_path
  providers = {
    vault      = vault
    kubernetes = kubernetes
  }
}

module "helm" {
  source           = "./modules/helm"
  namespace        = var.namespace
  chart_version    = var.chart_version
  existing_secret  = var.credentials_secret_name
  oidc_issuer_url  = var.oidc_issuer_url
  login_url        = var.login_url
  redeem_url       = var.redeem_url
  oidc_jwks_url    = var.oidc_jwks_url
  redirect_url     = var.redirect_url
  cookie_domains   = var.cookie_domains
  whitelist_domain = var.whitelist_domain
  oidc_scope       = var.oidc_scope
  cookie_refresh   = var.cookie_refresh
  session_store    = var.session_store
  redis_url        = var.redis_url
  service_port     = var.service_port
  providers = {
    helm = helm
  }
  depends_on = [module.secret]
}

resource "kubernetes_labels" "oauth2_proxy_use_waypoint" {
  depends_on  = [module.helm]
  api_version = "v1"
  kind        = "Service"
  metadata {
    name      = var.service_name
    namespace = var.namespace
  }
  labels = {
    "istio.io/use-waypoint"         = "bookstore-waypoint"
    "istio.io/ingress-use-waypoint" = "true"
  }
}

resource "kubernetes_manifest" "oauth2_proxy_timeout_route" {
  depends_on = [module.helm]
  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "HTTPRoute"
    metadata = {
      name      = "oauth2-proxy-timeout"
      namespace = var.namespace
    }
    spec = {
      parentRefs = [
        {
          kind  = "Service"
          group = ""
          name  = var.service_name
          port  = var.service_port
        }
      ]
      rules = [
        {
          backendRefs = [
            {
              name = var.service_name
              port = var.service_port
            }
          ]
          timeouts = {
            request = "15s"
          }
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "oauth2_proxy_authz" {
  depends_on = [module.helm]
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "oauth2-proxy-authz"
      namespace = var.namespace
    }
    spec = {
      selector = {
        matchLabels = {
          app = var.service_name
        }
      }
      action = "ALLOW"
      rules = [
        {
          from = [
            {
              source = {
                principals = [
                  "cluster.local/ns/istio-system/sa/bookstore-gateway-istio",
                  "cluster.local/ns/${var.namespace}/sa/bookstore-waypoint",
                ]
              }
            }
          ]
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "oauth2_proxy_route" {
  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "HTTPRoute"
    metadata = {
      name      = "oauth2-proxy-route"
      namespace = var.namespace
    }
    spec = {
      parentRefs = concat([
        {
          name        = var.gateway_name
          namespace   = var.gateway_namespace
          sectionName = "https-${split(".", var.app_host)[0]}"
        },
        {
          name        = var.gateway_name
          namespace   = var.gateway_namespace
          sectionName = "https-${split(".", var.admin_app_host)[0]}"
        }
        ], [
        for host in var.protected_hosts : {
          name        = var.gateway_name
          namespace   = var.gateway_namespace
          sectionName = "https-${split(".", host)[0]}"
        }
      ])
      hostnames = concat([var.app_host, var.admin_app_host], var.protected_hosts)
      rules = [
        {
          matches = [
            {
              path = {
                type  = "PathPrefix"
                value = "/oauth2/"
              }
            }
          ]
          backendRefs = [
            {
              group = ""
              kind  = "Service"
              name  = var.service_name
              port  = var.service_port
            }
          ]
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "admin_web_ext_authz" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "admin-web-edge-ext-authz"
      namespace = var.namespace
    }
    spec = {
      targetRefs = [
        {
          kind  = "Service"
          group = ""
          name  = var.admin_web_service_name
        }
      ]
      action = "CUSTOM"
      provider = {
        name = var.extension_provider_name
      }
      rules = [
        {
          to = [
            {
              operation = {
                notPaths = ["/oauth2/*"]
              }
            }
          ]
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "protected_host_ext_authz" {
  for_each = toset(var.protected_hosts)
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "${split(".", each.key)[0]}-edge-ext-authz"
      namespace = var.gateway_namespace
    }
    spec = {
      targetRefs = [
        {
          kind  = "Gateway"
          group = "gateway.networking.k8s.io"
          name  = var.gateway_name
        }
      ]
      action = "CUSTOM"
      provider = {
        name = var.extension_provider_name
      }
      rules = [
        {
          to = [
            {
              operation = {
                hosts    = [each.key, "${each.key}:*"]
                notPaths = ["/oauth2/*"]
              }
            }
          ]
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "admin_api_edge_ext_authz" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "admin-api-edge-ext-authz"
      namespace = var.gateway_namespace
    }
    spec = {
      targetRefs = [
        {
          kind  = "Gateway"
          group = "gateway.networking.k8s.io"
          name  = var.gateway_name
        }
      ]
      action = "CUSTOM"
      provider = {
        name = var.extension_provider_name
      }
      rules = [
        {
          to = [
            {
              operation = {
                hosts = [var.admin_app_host, "${var.admin_app_host}:*"]
                paths = ["/api/v1/admin", "/api/v1/admin/*"]
              }
            }
          ]
        }
      ]
    }
  }
}

