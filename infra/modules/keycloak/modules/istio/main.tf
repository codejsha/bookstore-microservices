resource "kubernetes_manifest" "istio_gateway" {
  manifest = {
    apiVersion = "networking.istio.io/v1"
    kind       = "Gateway"
    metadata = {
      name      = "keycloak-gw"
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
            var.address
          ]
        },
        {
          port = {
            number   = 443
            name     = "https"
            protocol = "HTTPS"
          }
          hosts = [
            var.address
          ]
          tls = {
            mode           = "SIMPLE"
            credentialName = "keycloak-cert"
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
      name      = "keycloak-vs"
      namespace = var.namespace
    }
    spec = {
      hosts = [
        var.address
      ]
      gateways = [
        kubernetes_manifest.istio_gateway.manifest.metadata.name
      ]
      http = [
        {
          route = [
            {
              destination = {
                host = "keycloak.keycloak.svc.cluster.local"
                port = {
                  number = 80
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
      name      = "keycloak-dr"
      namespace = var.namespace
    }
    spec = {
      host = "keycloak.keycloak.svc.cluster.local"
      trafficPolicy = {
        loadBalancer = {
          simple = "ROUND_ROBIN"
        }
      }
    }
  }
}
