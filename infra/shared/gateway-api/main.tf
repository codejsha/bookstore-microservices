locals {
  prefix = var.name_prefix != "" ? var.name_prefix : split(".", var.hostname)[0]
}

resource "kubernetes_manifest" "httproute" {
  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "HTTPRoute"
    metadata = {
      name      = "${local.prefix}-route"
      namespace = var.route_namespace
    }
    spec = {
      parentRefs = [
        {
          name        = var.gateway_name
          namespace   = var.gateway_namespace
          sectionName = "https-${split(".", var.hostname)[0]}"
        }
      ]
      hostnames = [var.hostname]
      rules = [
        merge(
          length(var.remove_request_headers) > 0 ? {
            filters = [
              {
                type                  = "RequestHeaderModifier"
                requestHeaderModifier = { remove = var.remove_request_headers }
              }
            ]
          } : {},
          {
            backendRefs = [
              {
                name      = var.service_name
                namespace = var.route_namespace
                port      = var.service_port
              }
            ]
          },
          var.request_timeout != null ? {
            timeouts = {
              request = var.request_timeout
            }
          } : {}
        )
      ]
    }
  }
}
