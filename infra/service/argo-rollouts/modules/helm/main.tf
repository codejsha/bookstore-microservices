terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "argo_rollouts" {
  namespace  = var.namespace
  name       = "argo-rollouts"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-rollouts"
  version    = "2.41.1"
  values = [
    file("${path.module}/values.yaml"),
    yamlencode({
      controller = {
        trafficRouterPlugins = [
          {
            name     = "argoproj-labs/gatewayAPI"
            location = var.gatewayapi_plugin_url
          },
        ]
      }
    }),
  ]
  timeout = 300
}
