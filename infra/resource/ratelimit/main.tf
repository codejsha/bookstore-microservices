terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
  }
}

locals {
  app_labels = {
    "app.kubernetes.io/name" = "ratelimit"
  }

  ratelimit_domain = "bookstore"

  ratelimit_config = yamlencode({
    domain = local.ratelimit_domain
    descriptors = [
      {
        key = "user_uid"
        rate_limit = {
          unit              = "second"
          requests_per_unit = var.user_requests_per_second
        }
      }
    ]
  })
}

resource "kubernetes_config_map_v1" "ratelimit_config" {
  metadata {
    name      = "ratelimit-config"
    namespace = var.namespace
    labels    = local.app_labels
  }
  data = {
    "bookstore.yaml" = local.ratelimit_config
  }
}

resource "kubernetes_deployment_v1" "ratelimit" {
  metadata {
    name      = "ratelimit"
    namespace = var.namespace
    labels    = local.app_labels
  }
  spec {
    replicas = var.replicas
    selector {
      match_labels = local.app_labels
    }
    template {
      metadata {
        labels = local.app_labels
        annotations = {
          "checksum/config" = sha256(local.ratelimit_config)
        }
      }
      spec {
        container {
          name  = "ratelimit"
          image = var.image
          args  = ["/bin/ratelimit"]

          env {
            name  = "RUNTIME_ROOT"
            value = "/data"
          }
          env {
            name  = "RUNTIME_SUBDIRECTORY"
            value = "ratelimit"
          }
          env {
            name  = "RUNTIME_WATCH_ROOT"
            value = "false"
          }
          env {
            name  = "RUNTIME_IGNOREDOTFILES"
            value = "true"
          }
          env {
            name  = "REDIS_SOCKET_TYPE"
            value = "tcp"
          }
          env {
            name  = "REDIS_URL"
            value = var.valkey_host
          }
          env {
            name  = "USE_STATSD"
            value = "false"
          }
          env {
            name  = "LOG_FORMAT"
            value = "json"
          }

          port {
            name           = "grpc"
            container_port = 8081
          }
          port {
            name           = "http"
            container_port = 8080
          }

          volume_mount {
            name       = "config"
            mount_path = "/data/ratelimit/config"
          }

          liveness_probe {
            http_get {
              path = "/healthcheck"
              port = 8080
            }
            initial_delay_seconds = 5
            period_seconds        = 10
          }

          resources {
            requests = {
              cpu    = "50m"
              memory = "64Mi"
            }
            limits = {
              cpu    = "500m"
              memory = "128Mi"
            }
          }
        }

        volume {
          name = "config"
          config_map {
            name = kubernetes_config_map_v1.ratelimit_config.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service_v1" "ratelimit" {
  metadata {
    name      = "ratelimit"
    namespace = var.namespace
    labels    = local.app_labels
  }
  spec {
    selector = local.app_labels
    port {
      name        = "grpc"
      port        = 8081
      target_port = 8081
    }
    port {
      name        = "http"
      port        = 8080
      target_port = 8080
    }
  }
}

resource "kubernetes_manifest" "waypoint_ratelimit_filter" {
  manifest = {
    apiVersion = "networking.istio.io/v1alpha3"
    kind       = "EnvoyFilter"
    metadata = {
      name      = "waypoint-user-ratelimit"
      namespace = var.namespace
    }
    spec = {
      priority = 10
      workloadSelector = {
        labels = {
          "gateway.networking.k8s.io/gateway-name" = var.waypoint_name
        }
      }
      configPatches = [
        {
          applyTo = "CLUSTER"
          match = {
            context = "GATEWAY"
          }
          patch = {
            operation = "ADD"
            value = {
              name            = "ratelimit_service"
              type            = "STRICT_DNS"
              connect_timeout = "1s"
              lb_policy       = "ROUND_ROBIN"
              typed_extension_protocol_options = {
                "envoy.extensions.upstreams.http.v3.HttpProtocolOptions" = {
                  "@type" = "type.googleapis.com/envoy.extensions.upstreams.http.v3.HttpProtocolOptions"
                  explicit_http_config = {
                    http2_protocol_options = {}
                  }
                }
              }
              load_assignment = {
                cluster_name = "ratelimit_service"
                endpoints = [
                  {
                    lb_endpoints = [
                      {
                        endpoint = {
                          address = {
                            socket_address = {
                              address    = "ratelimit.${var.namespace}.svc.cluster.local"
                              port_value = 8081
                            }
                          }
                        }
                      }
                    ]
                  }
                ]
              }
            }
          }
        },
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
              name = "envoy.filters.http.ratelimit"
              typed_config = {
                "@type"           = "type.googleapis.com/envoy.extensions.filters.http.ratelimit.v3.RateLimit"
                domain            = local.ratelimit_domain
                failure_mode_deny = false
                timeout           = "0.25s"
                rate_limit_service = {
                  transport_api_version = "V3"
                  grpc_service = {
                    envoy_grpc = {
                      cluster_name = "ratelimit_service"
                    }
                  }
                }
              }
            }
          }
        },
        {
          applyTo = "VIRTUAL_HOST"
          match = {
            context = "GATEWAY"
          }
          patch = {
            operation = "MERGE"
            value = {
              rate_limits = [
                {
                  actions = [
                    {
                      request_headers = {
                        header_name    = "x-user-id"
                        descriptor_key = "user_uid"
                        skip_if_absent = true
                      }
                    }
                  ]
                }
              ]
            }
          }
        }
      ]
    }
  }
}
