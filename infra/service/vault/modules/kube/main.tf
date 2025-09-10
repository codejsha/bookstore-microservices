terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
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
  token_reviewer_jwt     = data.kubernetes_secret_v1.vault_token.data["token"]
  issuer                 = "https://kubernetes.default.svc.cluster.local"
  disable_iss_validation = true
  disable_local_ca_jwt   = true
}

resource "kubernetes_service_account_v1" "tf_auth" {
  for_each = var.tf_auth_roles
  metadata {
    name      = "tf-${each.key}"
    namespace = var.namespace
  }
}

resource "vault_kubernetes_auth_backend_role" "tf_auth" {
  for_each                         = var.tf_auth_roles
  backend                          = vault_auth_backend.kubernetes.path
  role_name                        = "tf-${each.key}"
  bound_service_account_names      = [kubernetes_service_account_v1.tf_auth[each.key].metadata[0].name]
  bound_service_account_namespaces = [var.namespace]
  token_policies                   = [each.value]
  token_ttl                        = 1200
}
