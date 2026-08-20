terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "vault_transit" {
  namespace  = var.namespace
  name       = "vault"
  repository = "https://helm.releases.hashicorp.com"
  chart      = "vault"
  version    = "0.32.0"
  values = [
    file("${path.module}/values.yaml")
  ]
  timeout = 120
}
