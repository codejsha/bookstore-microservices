terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
  }
}

resource "vault_policy" "opensearch" {
  name = "opensearch"
  policy = file("${path.module}/policy.hcl")
}

resource "vault_kubernetes_auth_backend_role" "opensearch" {
  role_name = "opensearch-role"
  bound_service_account_names = ["opensearch"]
  bound_service_account_namespaces = [var.namespace]
  token_policies = [vault_policy.opensearch.name]
  token_ttl = "3600" # 1 hour
}

resource "vault_kv_secret_v2" "opensearch" {
  name  = "opensearch/admin"
  mount = "kv"
  data_json = jsonencode({
    initial_password = var.initial_admin_password
  })
}
