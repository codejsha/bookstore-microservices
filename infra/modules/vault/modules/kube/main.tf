terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
  }
}

resource "vault_auth_backend" "kubernetes" {
  type = "kubernetes"
}

resource "vault_kubernetes_auth_backend_config" "kubernetes" {
  backend                = vault_auth_backend.kubernetes.path
  kubernetes_host        = "https://${var.kube_api_server_address}"
  kubernetes_ca_cert     = var.kube_ca_crt
  token_reviewer_jwt     = "$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)"
  issuer                 = "https://kubernetes.default.svc.cluster.local"
  disable_iss_validation = true
  disable_local_ca_jwt   = true
}
