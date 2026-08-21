terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "reflector" {
  namespace  = var.namespace
  name       = "reflector"
  repository = "oci://ghcr.io/emberstack/helm-charts"
  chart      = "reflector"
  version    = "9.1.10"
}
