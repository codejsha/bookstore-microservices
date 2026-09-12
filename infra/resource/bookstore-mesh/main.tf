terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
  }
}

locals {
  service_account = {
    catalog      = "catalog-postgres"
    customer     = "customer-postgres"
    identity     = "identity-postgres"
    inventory    = "inventory-postgres"
    order        = "order-mysql"
    payment      = "payment-mysql"
    delivery     = "delivery-mysql"
    notification = "notification-mysql"
    support      = "support-mysql"
    settlement   = "settlement-mysql"
    admin        = "admin"
    admin_web    = "admin-web"
  }
  principal = {
    for svc, sa in local.service_account : svc => "cluster.local/ns/${var.namespace}/sa/${sa}"
  }

  authz_check_paths = ["/internal/authz", "/internal/authz/*"]

  webhook_paths = ["/internal/webhooks/*"]

  unauthenticated_paths = concat(var.anonymous_paths, local.authz_check_paths, local.webhook_paths)

  waypoint_name      = "bookstore-waypoint"
  waypoint_principal = "cluster.local/ns/${var.namespace}/sa/${local.waypoint_name}"

  role_staff = "STAFF"

  introspect_services = [
    "admin",
    "catalog",
    "customer",
    "delivery",
    "identity",
    "inventory",
    "notification",
    "order",
    "payment",
    "settlement",
    "support",
  ]

  timeout_services = [
    "admin",
    "admin-web",
    "settlement",
    "web",
  ]

  circuit_breaker_services = [
    "catalog",
    "identity",
    "order",
    "payment",
  ]

  service_port = 8080
}

resource "kubernetes_manifest" "peer_auth_strict" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "PeerAuthentication"
    metadata = {
      name      = "default"
      namespace = var.namespace
    }
    spec = {
      mtls = {
        mode = "STRICT"
      }
    }
  }
}

resource "kubernetes_manifest" "request_auth_keycloak" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "RequestAuthentication"
    metadata = {
      name      = "bookstore-keycloak-jwt"
      namespace = var.namespace
    }
    spec = {
      targetRefs = [
        {
          group = "gateway.networking.k8s.io"
          kind  = "Gateway"
          name  = local.waypoint_name
        }
      ]
      jwtRules = [
        merge(
          {
            issuer                = var.keycloak_issuer
            jwksUri               = var.keycloak_jwks_uri
            forwardOriginalToken  = true
            outputPayloadToHeader = "x-jwt-payload"
            outputClaimToHeaders = [
              { header = "X-User-Id", claim = "sub" },
              { header = "X-User-Email", claim = "email" },
              { header = "X-User-Name", claim = "preferred_username" },
              { header = "X-User-Roles", claim = "roles" },
              { header = "X-User-Scopes", claim = "scope" },
            ]
          },
          length(var.keycloak_audiences) > 0 ? { audiences = var.keycloak_audiences } : {}
        )
      ]
    }
  }
}

resource "kubernetes_manifest" "require_jwt" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "require-jwt"
      namespace = var.namespace
    }
    spec = {
      action = "DENY"
      rules = concat(
        [
          {
            from = [
              { source = { notRequestPrincipals = ["*"] } }
            ]
            to = [
              {
                operation = {
                  notPaths = concat(local.unauthenticated_paths, var.anonymous_read_paths)
                  notPorts = ["9090"]
                }
              }
            ]
          }
        ],
        length(var.anonymous_read_paths) == 0 ? [] : [
          {
            from = [
              { source = { notRequestPrincipals = ["*"] } }
            ]
            to = [
              {
                operation = {
                  paths      = var.anonymous_read_paths
                  notMethods = ["GET"]
                  notPorts   = ["9090"]
                }
              }
            ]
          }
        ]
      )
    }
  }
}

resource "kubernetes_manifest" "l4_waypoint_only" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "bookstore-l4-waypoint-only"
      namespace = var.namespace
    }
    spec = {
      selector = {
        matchLabels = {
          "app.kubernetes.io/group" = "bookstore"
        }
      }
      action = "ALLOW"
      rules = [
        {
          from = [
            { source = { principals = [local.waypoint_principal] } }
          ]
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "waypoint" {
  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "Gateway"
    metadata = {
      name      = local.waypoint_name
      namespace = var.namespace
      labels = {
        "istio.io/waypoint-for" = "service"
      }
    }
    spec = {
      gatewayClassName = "istio-waypoint"
      listeners = [
        {
          name     = "mesh"
          port     = 15008
          protocol = "HBONE"
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "request_timeout" {
  for_each = toset(local.timeout_services)

  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "HTTPRoute"
    metadata = {
      name      = "${each.key}-timeout"
      namespace = var.namespace
    }
    spec = {
      parentRefs = [
        {
          kind  = "Service"
          group = ""
          name  = each.key
          port  = local.service_port
        }
      ]
      rules = [
        {
          backendRefs = [
            {
              name = each.key
              port = local.service_port
            }
          ]
          timeouts = {
            request = var.request_timeout
          }
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "authz_identity" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "identity-authz"
      namespace = var.namespace
    }
    spec = {
      targetRefs = [
        {
          kind  = "Service"
          group = ""
          name  = "identity"
        }
      ]
      action = "ALLOW"
      rules = [
        {
          from = [
            {
              source = {
                principals = [
                  local.principal.admin,
                  local.principal.customer,
                  var.ingress_gateway_principal,
                ]
              }
            }
          ]
          to = [
            { operation = { notPaths = local.authz_check_paths } }
          ]
        },
        {
          from = [
            { source = { principals = [local.waypoint_principal] } }
          ]
          to = [
            { operation = { paths = local.authz_check_paths } }
          ]
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "introspect_authz" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "bookstore-introspect"
      namespace = var.namespace
    }
    spec = {
      targetRefs = [
        for svc in local.introspect_services : {
          kind  = "Service"
          group = ""
          name  = svc
        }
      ]
      action = "CUSTOM"
      provider = {
        name = "keycloak-introspect"
      }
      rules = concat(
        [
          {
            to = [
              {
                operation = {
                  notPaths = concat(local.unauthenticated_paths, var.anonymous_read_paths)
                  notPorts = ["9090"]
                }
              }
            ]
          }
        ],
        length(var.anonymous_read_paths) == 0 ? [] : [
          {
            to = [
              {
                operation = {
                  paths      = var.anonymous_read_paths
                  notMethods = ["GET"]
                  notPorts   = ["9090"]
                }
              }
            ]
          }
        ]
      )
    }
  }
}

resource "kubernetes_manifest" "authz_admin" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "admin-authz"
      namespace = var.namespace
    }
    spec = {
      targetRefs = [
        {
          kind  = "Service"
          group = ""
          name  = "admin"
        }
      ]
      action = "ALLOW"
      rules = [
        {
          from = [
            {
              source = {
                principals = [
                  local.principal.admin_web,
                  var.ingress_gateway_principal,
                ]
              }
            }
          ]
          when = [
            {
              key    = "request.auth.claims[roles]"
              values = [local.role_staff]
            }
          ]
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "authz_catalog" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "catalog-authz"
      namespace = var.namespace
    }
    spec = {
      targetRefs = [
        {
          kind  = "Service"
          group = ""
          name  = "catalog"
        }
      ]
      action = "ALLOW"
      rules = [
        {
          from = [
            {
              source = {
                principals = [
                  local.principal.admin,
                  local.principal.order,
                  local.principal.customer,
                  var.ingress_gateway_principal,
                ]
              }
            }
          ]
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "authz_customer" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "customer-authz"
      namespace = var.namespace
    }
    spec = {
      targetRefs = [
        {
          kind  = "Service"
          group = ""
          name  = "customer"
        }
      ]
      action = "ALLOW"
      rules = [
        {
          from = [
            {
              source = {
                principals = [
                  local.principal.admin,
                  var.ingress_gateway_principal,
                ]
              }
            }
          ]
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "authz_inventory" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "inventory-authz"
      namespace = var.namespace
    }
    spec = {
      targetRefs = [
        {
          kind  = "Service"
          group = ""
          name  = "inventory"
        }
      ]
      action = "ALLOW"
      rules = [
        {
          from = [
            {
              source = {
                principals = [
                  local.principal.admin,
                  local.principal.catalog,
                  local.principal.order,
                  var.ingress_gateway_principal,
                ]
              }
            }
          ]
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "authz_order" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "order-authz"
      namespace = var.namespace
    }
    spec = {
      targetRefs = [
        {
          kind  = "Service"
          group = ""
          name  = "order"
        }
      ]
      action = "ALLOW"
      rules = [
        {
          from = [
            {
              source = {
                principals = [
                  local.principal.admin,
                  local.principal.customer,
                  local.principal.payment,
                  var.ingress_gateway_principal,
                ]
              }
            }
          ]
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "authz_payment" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "payment-authz"
      namespace = var.namespace
    }
    spec = {
      targetRefs = [
        {
          kind  = "Service"
          group = ""
          name  = "payment"
        }
      ]
      action = "ALLOW"
      rules = [
        {
          from = [
            {
              source = {
                principals = [
                  local.principal.admin,
                  local.principal.customer,
                  local.principal.order,
                  var.ingress_gateway_principal,
                ]
              }
            }
          ]
          to = [
            { operation = { notPaths = local.webhook_paths } }
          ]
        },
        {
          from = [
            { source = { principals = [var.hyperswitch_router_principal] } }
          ]
          to = [
            { operation = { paths = local.webhook_paths, methods = ["POST"] } }
          ]
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "authz_support" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "support-authz"
      namespace = var.namespace
    }
    spec = {
      targetRefs = [
        {
          kind  = "Service"
          group = ""
          name  = "support"
        }
      ]
      action = "ALLOW"
      rules = [
        {
          from = [
            {
              source = {
                principals = [
                  local.principal.customer,
                  var.ingress_gateway_principal,
                ]
              }
            }
          ]
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "authz_settlement" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "settlement-authz"
      namespace = var.namespace
    }
    spec = {
      targetRefs = [
        {
          kind  = "Service"
          group = ""
          name  = "settlement"
        }
      ]
      action = "ALLOW"
      rules = [
        {
          from = [
            {
              source = {
                principals = [
                  local.principal.admin,
                  var.ingress_gateway_principal,
                ]
              }
            }
          ]
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "authz_delivery" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "delivery-authz"
      namespace = var.namespace
    }
    spec = {
      targetRefs = [
        {
          kind  = "Service"
          group = ""
          name  = "delivery"
        }
      ]
      action = "ALLOW"
      rules = [
        {
          from = [
            {
              source = {
                principals = [
                  local.principal.order,
                  local.principal.customer,
                  var.ingress_gateway_principal,
                ]
              }
            }
          ]
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "circuit_breaker" {
  for_each = toset(local.circuit_breaker_services)

  manifest = {
    apiVersion = "networking.istio.io/v1"
    kind       = "DestinationRule"
    metadata = {
      name      = "${each.key}-circuit-breaker"
      namespace = var.namespace
    }
    spec = {
      host = "${each.key}.${var.namespace}.svc.cluster.local"
      trafficPolicy = {
        connectionPool = {
          http = {
            http2MaxRequests         = 256
            maxRequestsPerConnection = 0
          }
        }
        outlierDetection = {
          consecutive5xxErrors = 5
          interval             = "10s"
          baseEjectionTime     = "30s"
          maxEjectionPercent   = 50
        }
      }
    }
  }
}

resource "kubernetes_manifest" "authz_notification" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "notification-authz"
      namespace = var.namespace
    }
    spec = {
      targetRefs = [
        {
          kind  = "Service"
          group = ""
          name  = "notification"
        }
      ]
      action = "ALLOW"
      rules = [
        {
          from = [
            {
              source = {
                principals = [
                  var.ingress_gateway_principal,
                ]
              }
            }
          ]
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "mesh_tracing" {
  manifest = {
    apiVersion = "telemetry.istio.io/v1"
    kind       = "Telemetry"
    metadata = {
      name      = "mesh-tracing"
      namespace = var.namespace
    }
    spec = {
      tracing = [
        {
          providers = [
            { name = "otel-tracing" }
          ]
          randomSamplingPercentage = 10
        }
      ]
    }
  }
}
