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
  }
}

resource "kubernetes_namespace_v1" "istio-system" {
  metadata {
    name = var.namespace
  }
}

module "helm" {
  source        = "./modules/helm"
  namespace     = kubernetes_namespace_v1.istio-system.metadata[0].name
  istio_version = var.istio_version
  providers = {
    helm = helm
  }
}

module "kiali" {
  source              = "./modules/kiali"
  namespace           = kubernetes_namespace_v1.istio-system.metadata[0].name
  kiali_chart_version = var.kiali_chart_version
  grafana_username    = var.grafana_username
  grafana_password    = var.grafana_password
  providers = {
    helm = helm
  }
}

data "kubernetes_config_map_v1" "kube_root_ca" {
  metadata {
    name      = "kube-root-ca.crt"
    namespace = "kube-system"
  }
}

resource "kubernetes_service_account_v1" "wildcard_issuer" {
  metadata {
    name      = "wildcard-issuer"
    namespace = kubernetes_namespace_v1.istio-system.metadata[0].name
  }
}

resource "kubernetes_secret_v1" "wildcard_issuer_token" {
  type = "kubernetes.io/service-account-token"
  metadata {
    name      = "wildcard-issuer-token"
    namespace = kubernetes_namespace_v1.istio-system.metadata[0].name
    annotations = {
      "kubernetes.io/service-account.name" = kubernetes_service_account_v1.wildcard_issuer.metadata[0].name
    }
  }
}

resource "kubernetes_cluster_role_binding_v1" "wildcard_issuer_token_rolebinding" {
  metadata {
    name = "wildcard-issuer-token-rolebinding"
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = "system:auth-delegator"
  }
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account_v1.wildcard_issuer.metadata[0].name
    namespace = kubernetes_namespace_v1.istio-system.metadata[0].name
  }
}

resource "kubernetes_manifest" "wildcard_issuer" {
  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "Issuer"
    metadata = {
      name      = "wildcard-issuer"
      namespace = kubernetes_namespace_v1.istio-system.metadata[0].name
    }
    spec = {
      vault = {
        server   = "http://vault.vault.svc.cluster.local:8200"
        path     = "pki_int/sign/wildcard"
        caBundle = base64encode(trimspace(data.kubernetes_config_map_v1.kube_root_ca.data["ca.crt"]))
        auth = {
          kubernetes = {
            mountPath = "/v1/auth/kubernetes"
            role      = "wildcard-issuer"
            secretRef = {
              name = kubernetes_secret_v1.wildcard_issuer_token.metadata[0].name
              key  = "token"
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_manifest" "wildcard_cert" {
  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "Certificate"
    metadata = {
      name      = "wildcard-cert"
      namespace = kubernetes_namespace_v1.istio-system.metadata[0].name
    }
    spec = {
      secretName = "wildcard-cert"
      commonName = "example.com"
      dnsNames = [
        "*.example.com",
        "example.com",
      ]
      issuerRef = {
        group = "cert-manager.io"
        kind  = "Issuer"
        name  = kubernetes_manifest.wildcard_issuer.manifest.metadata.name
      }
    }
  }
}

resource "kubernetes_manifest" "edge_host_cert" {
  for_each = var.edge_hosts
  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "Certificate"
    metadata = {
      name      = "${split(".", each.key)[0]}-edge-cert"
      namespace = kubernetes_namespace_v1.istio-system.metadata[0].name
    }
    spec = {
      secretName = "${split(".", each.key)[0]}-edge-cert"
      commonName = each.key
      dnsNames   = [each.key]
      issuerRef = {
        group = "cert-manager.io"
        kind  = "Issuer"
        name  = kubernetes_manifest.wildcard_issuer.manifest.metadata.name
      }
    }
  }
}

locals {
  edge_host_listeners = [
    for host in sort(keys(var.edge_hosts)) : {
      name     = "https-${split(".", host)[0]}"
      port     = 443
      protocol = "HTTPS"
      hostname = host
      tls = {
        mode = "Terminate"
        certificateRefs = [
          { name = kubernetes_manifest.edge_host_cert[host].manifest.spec.secretName }
        ]
      }
      allowedRoutes = {
        namespaces = {
          from = "Selector"
          selector = {
            matchExpressions = [
              {
                key      = "kubernetes.io/metadata.name"
                operator = "In"
                values   = var.edge_hosts[host]
              }
            ]
          }
        }
      }
    }
  ]
}

resource "kubernetes_config_map_v1" "bookstore_gateway_options" {
  metadata {
    name      = "bookstore-gateway-options"
    namespace = kubernetes_namespace_v1.istio-system.metadata[0].name
  }
  data = {
    service = yamlencode({
      spec = {
        externalTrafficPolicy = "Local"
      }
    })
  }
}

resource "kubernetes_manifest" "bookstore_gateway" {
  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "Gateway"
    metadata = {
      name      = "bookstore-gateway"
      namespace = kubernetes_namespace_v1.istio-system.metadata[0].name
    }
    spec = {
      gatewayClassName = "istio"
      infrastructure = {
        parametersRef = {
          group = ""
          kind  = "ConfigMap"
          name  = kubernetes_config_map_v1.bookstore_gateway_options.metadata[0].name
        }
      }
      listeners = concat([
        {
          name     = "http"
          port     = 80
          protocol = "HTTP"
          hostname = "*.example.com"
          allowedRoutes = {
            namespaces = {
              from = "Same"
            }
          }
        }
      ], local.edge_host_listeners)
    }
  }
}

resource "kubernetes_manifest" "edge_blocklist" {
  count = length(var.edge_blocked_cidrs) > 0 ? 1 : 0

  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "edge-blocklist"
      namespace = kubernetes_namespace_v1.istio-system.metadata[0].name
    }
    spec = {
      targetRefs = [
        {
          kind  = "Gateway"
          group = "gateway.networking.k8s.io"
          name  = kubernetes_manifest.bookstore_gateway.manifest.metadata.name
        }
      ]
      action = "DENY"
      rules = [
        {
          from = [
            { source = { remoteIpBlocks = var.edge_blocked_cidrs } }
          ]
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "http_to_https_redirect" {
  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "HTTPRoute"
    metadata = {
      name      = "http-to-https-redirect"
      namespace = kubernetes_namespace_v1.istio-system.metadata[0].name
    }
    spec = {
      parentRefs = [
        {
          name        = kubernetes_manifest.bookstore_gateway.manifest.metadata.name
          namespace   = kubernetes_namespace_v1.istio-system.metadata[0].name
          sectionName = "http"
        }
      ]
      rules = [
        {
          filters = [
            {
              type = "RequestRedirect"
              requestRedirect = {
                scheme     = "https"
                statusCode = 301
              }
            }
          ]
        }
      ]
    }
  }
}

module "kiali_route" {
  source            = "../../shared/gateway-api"
  hostname          = var.kiali_address
  service_name      = var.kiali_service_name
  service_port      = 20001
  route_namespace   = kubernetes_namespace_v1.istio-system.metadata[0].name
  gateway_name      = kubernetes_manifest.bookstore_gateway.manifest.metadata.name
  gateway_namespace = kubernetes_namespace_v1.istio-system.metadata[0].name
  request_timeout   = "10m"
}

locals {
  edge_ratelimit_filter_type = "type.googleapis.com/envoy.extensions.filters.http.local_ratelimit.v3.LocalRateLimit"

  edge_ratelimit_hosts = {
    "keycloak.example.com:443" = {
      stat_prefix   = "keycloak_edge_ratelimit"
      host_bucket   = { max_tokens = 200, tokens_per_fill = 100, fill_interval = "1s" }
      per_ip_bucket = { max_tokens = 30, tokens_per_fill = 15, fill_interval = "1s" }
      paths = {
        auth_token = {
          match         = { prefix_match = "/realms/${var.keycloak_realm}/protocol/openid-connect/token" }
          bucket        = { max_tokens = 40, tokens_per_fill = 20, fill_interval = "1s" }
          per_ip_bucket = { max_tokens = 10, tokens_per_fill = 5, fill_interval = "1s" }
        }
        auth_login = {
          match         = { prefix_match = "/realms/${var.keycloak_realm}/login-actions/" }
          bucket        = { max_tokens = 20, tokens_per_fill = 10, fill_interval = "1s" }
          per_ip_bucket = { max_tokens = 6, tokens_per_fill = 3, fill_interval = "1s" }
        }
      }
    }
    "identity-api.example.com:443" = {
      stat_prefix   = "identity_edge_ratelimit"
      host_bucket   = { max_tokens = 50, tokens_per_fill = 25, fill_interval = "1s" }
      per_ip_bucket = { max_tokens = 20, tokens_per_fill = 10, fill_interval = "1s" }
      paths         = {}
    }
    "order-api.example.com:443" = {
      stat_prefix   = "order_edge_ratelimit"
      host_bucket   = { max_tokens = 200, tokens_per_fill = 100, fill_interval = "1s" }
      per_ip_bucket = { max_tokens = 40, tokens_per_fill = 20, fill_interval = "1s" }
      paths = {
        order_saga = {
          match         = { safe_regex_match = { regex = "^/api/v1/orders/(saga/place|[^/]+/(cancellation|fulfill))$" } }
          bucket        = { max_tokens = 20, tokens_per_fill = 10, fill_interval = "1s" }
          per_ip_bucket = { max_tokens = 10, tokens_per_fill = 5, fill_interval = "1s" }
        }
      }
    }
    "catalog-api.example.com:443" = {
      stat_prefix   = "catalog_edge_ratelimit"
      host_bucket   = { max_tokens = 1000, tokens_per_fill = 500, fill_interval = "1s" }
      per_ip_bucket = { max_tokens = 60, tokens_per_fill = 30, fill_interval = "1s" }
      paths         = {}
    }
  }

  edge_ratelimit_per_ip = var.edge_per_ip_ratelimit_enabled

  edge_ratelimit_vhosts = {
    for vhost, cfg in local.edge_ratelimit_hosts : vhost => {
      stat_prefix = cfg.stat_prefix
      host_bucket = cfg.host_bucket
      rate_limits = concat(
        [
          for key, path in cfg.paths : {
            actions = [
              {
                header_value_match = {
                  descriptor_value = key
                  headers          = [merge({ name = ":path" }, path.match)]
                }
              }
            ]
          }
        ],
        local.edge_ratelimit_per_ip && try(cfg.per_ip_bucket, null) != null ? [
          { actions = [{ remote_address = {} }] }
        ] : [],
        local.edge_ratelimit_per_ip ? [
          for key, path in cfg.paths : {
            actions = [
              {
                header_value_match = {
                  descriptor_value = key
                  headers          = [merge({ name = ":path" }, path.match)]
                }
              },
              { remote_address = {} },
            ]
          } if try(path.per_ip_bucket, null) != null
        ] : []
      )
      descriptors = concat(
        [
          for key, path in cfg.paths : {
            entries      = [{ key = "header_match", value = key }]
            token_bucket = path.bucket
          }
        ],
        local.edge_ratelimit_per_ip && try(cfg.per_ip_bucket, null) != null ? [
          {
            entries      = [{ key = "remote_address", value = "" }]
            token_bucket = cfg.per_ip_bucket
          }
        ] : [],
        local.edge_ratelimit_per_ip ? [
          for key, path in cfg.paths : {
            entries = [
              { key = "header_match", value = key },
              { key = "remote_address", value = "" },
            ]
            token_bucket = path.per_ip_bucket
          } if try(path.per_ip_bucket, null) != null
        ] : []
      )
    }
  }
}

resource "kubernetes_manifest" "edge_ratelimit" {
  manifest = {
    apiVersion = "networking.istio.io/v1alpha3"
    kind       = "EnvoyFilter"
    metadata = {
      name      = "edge-ratelimit"
      namespace = kubernetes_namespace_v1.istio-system.metadata[0].name
    }
    spec = {
      priority = 10
      workloadSelector = {
        labels = {
          "gateway.networking.k8s.io/gateway-name" = "bookstore-gateway"
        }
      }
      configPatches = concat(
        [
          {
            applyTo = "HTTP_FILTER"
            match = {
              context = "GATEWAY"
              listener = {
                filterChain = {
                  filter = {
                    name = "envoy.filters.network.http_connection_manager"
                    subFilter = {
                      name = "envoy.filters.http.router"
                    }
                  }
                }
              }
            }
            patch = {
              operation = "INSERT_BEFORE"
              value = {
                name = "envoy.filters.http.local_ratelimit"
                typed_config = {
                  "@type"  = "type.googleapis.com/udpa.type.v1.TypedStruct"
                  type_url = local.edge_ratelimit_filter_type
                  value = {
                    stat_prefix = "edge_local_ratelimit"
                  }
                }
              }
            }
          }
        ],
        [
          for vhost, cfg in local.edge_ratelimit_vhosts : {
            applyTo = "VIRTUAL_HOST"
            match = {
              context = "GATEWAY"
              routeConfiguration = {
                vhost = {
                  name = vhost
                }
              }
            }
            patch = {
              operation = "MERGE"
              value = {
                rate_limits = cfg.rate_limits
                typed_per_filter_config = {
                  "envoy.filters.http.local_ratelimit" = {
                    "@type"  = "type.googleapis.com/udpa.type.v1.TypedStruct"
                    type_url = local.edge_ratelimit_filter_type
                    value = merge(
                      {
                        stat_prefix  = cfg.stat_prefix
                        token_bucket = cfg.host_bucket
                        descriptors  = cfg.descriptors
                        filter_enabled = {
                          runtime_key = "edge_ratelimit_enabled"
                          default_value = {
                            numerator   = 100
                            denominator = "HUNDRED"
                          }
                        }
                        filter_enforced = {
                          runtime_key = "edge_ratelimit_enforced"
                          default_value = {
                            numerator   = 100
                            denominator = "HUNDRED"
                          }
                        }
                        response_headers_to_add = [
                          {
                            header        = { key = "retry-after", value = "1" }
                            append_action = "OVERWRITE_IF_EXISTS_OR_ADD"
                          }
                        ]
                      },
                      local.edge_ratelimit_per_ip ? { max_dynamic_descriptors = var.edge_per_ip_max_tracked_addresses } : {}
                    )
                  }
                }
              }
            }
          }
        ]
      )
    }
  }
}
