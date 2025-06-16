terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
  }
}

resource "vault_policy" "gitea" {
  name = "gitea"
  policy = file("${path.module}/policy.hcl")
}

resource "vault_kubernetes_auth_backend_role" "gitea" {
  role_name = "gitea-role"
  bound_service_account_names = ["gitea"]
  bound_service_account_namespaces = [var.namespace]
  token_policies = [vault_policy.gitea.name]
  token_ttl = "3600" # 1 hour
}

resource "vault_kv_secret_v2" "gitea" {
  name  = "gitea/admin/credentials"
  mount = "kv"
  data_json = jsonencode({
    email    = "${var.admin_username}@example.com"
    username = var.admin_username,
    password = var.admin_password
  })
}
