terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "hyperswitch" {
  namespace  = var.namespace
  name       = "hyperswitch"
  repository = "https://juspay.github.io/hyperswitch-helm"
  chart      = "hyperswitch-stack"
  version    = "0.2.20"
  values = [
    file("${path.module}/values.yaml")
  ]
}
