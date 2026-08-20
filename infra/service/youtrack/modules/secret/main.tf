terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
  }
}

resource "vault_policy" "youtrack" {
  name   = "youtrack"
  policy = file("${path.module}/policy.hcl")
}

resource "vault_kubernetes_auth_backend_role" "youtrack" {
  role_name                        = "youtrack-role"
  bound_service_account_names      = ["youtrack"]
  bound_service_account_namespaces = [var.namespace]
  token_policies                   = [vault_policy.youtrack.name]
  token_ttl                        = "3600"
}

resource "vault_kv_secret_v2" "license" {
  name  = "youtrack/license/credentials"
  mount = "kv"
  data_json = jsonencode({
    license_name = ""
    license_key  = ""
    operator     = var.admin_email
  })
  lifecycle {
    ignore_changes = [data_json]
  }
}
