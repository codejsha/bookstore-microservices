resource "kubernetes_manifest" "istio_gateway" {
  manifest = {
    apiVersion = "networking.istio.io/v1"
    kind       = "Gateway"
    metadata = {
      name      = "${var.name_prefix}-gw"
      namespace = var.namespace
    }
    spec = {
      selector = {
        istio = "ingressgateway"
      }
      servers = [
        {
          port = {
            number   = 80
            name     = "http"
            protocol = "HTTP"
          }
          hosts = [
            var.host_address
          ]
        },
        {
          port = {
            number   = 443
            name     = "https"
            protocol = "HTTPS"
          }
          hosts = [
            var.host_address
          ]
          tls = {
            mode           = "SIMPLE"
            credentialName = "${var.name_prefix}-cert"
          }
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "istio_virtual_service" {
  manifest = {
    apiVersion = "networking.istio.io/v1"
    kind       = "VirtualService"
    metadata = {
      name      = "${var.name_prefix}-vs"
      namespace = var.namespace
    }
    spec = {
      hosts = [
        var.host_address
      ]
      gateways = [
        kubernetes_manifest.istio_gateway.manifest.metadata.name,
      ]
      http = [
        {
          match = [
            {
              uri = {
                prefix = "/v1/metrics"
              }
            }
          ]
          route = [
            {
              destination = {
                host = var.host_fqdn
                port = {
                  number = 4318
                }
              }
            }
          ]
        },
        {
          route = [
            {
              destination = {
                host = var.host_fqdn
                port = {
                  number = 4317
                }
              }
            }
          ]
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "istio_destination_rule" {
  manifest = {
    apiVersion = "networking.istio.io/v1"
    kind       = "DestinationRule"
    metadata = {
      name      = "${var.name_prefix}-dr"
      namespace = var.namespace
    }
    spec = {
      host = var.host_fqdn
      trafficPolicy = {
        loadBalancer = {
          simple = "ROUND_ROBIN"
        },
      }
    }
  }
}
