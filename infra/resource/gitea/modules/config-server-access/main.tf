terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
  }
}

resource "vault_policy" "config_server_gitea" {
  name   = "config-server-gitea"
  policy = file("${path.module}/policy.hcl")
}

resource "vault_kubernetes_auth_backend_role" "config_server_gitea" {
  role_name                        = "config-server-gitea-role"
  bound_service_account_names      = [var.service_account]
  bound_service_account_namespaces = [var.namespace]
  token_policies                   = [vault_policy.config_server_gitea.name]
  token_ttl                        = "3600"
}
