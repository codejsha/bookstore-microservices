terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "reflector" {
  namespace = var.namespace
  name      = "reflector"
  chart     = "oci://ghcr.io/emberstack/helm-charts/reflector"
  version   = "10.0.20"
  values = [
    yamlencode({
      resources = {
        requests = { cpu = "10m", memory = "64Mi" }
        limits   = { cpu = "200m", memory = "256Mi" }
      }
    })
  ]
}
