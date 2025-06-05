terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "vault" {
  namespace  = var.namespace
  name       = "vault"
  repository = "https://helm.releases.hashicorp.com"
  chart      = "vault"
  version    = "0.30.0"
  values = [
    file("${path.module}/values.yaml")
  ]
  timeout = 120
}

# resource "helm_release" "vault_bitnami" {
#   namespace = var.namespace
#   name      = "vault"
#   chart     = "${path.module}/vault"
#   values = [
#     file("${path.module}/values-bitnami.yaml")
#   ]
#   timeout = 120
# }
