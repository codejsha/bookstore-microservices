terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "vso" {
  namespace  = var.namespace
  name       = "vault-secrets-operator"
  repository = "https://helm.releases.hashicorp.com"
  chart      = "vault-secrets-operator"
  version    = "1.4.0"

  values = [
    yamlencode({
      defaultVaultConnection = {
        enabled = true
        address = var.vault_address
      }
    }),
  ]
}
