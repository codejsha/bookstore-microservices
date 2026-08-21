terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
  }
}

resource "vault_policy" "bookstore_mysql" {
  name = "bookstore-mysql"
  policy = file("${path.module}/policy.hcl")
}

resource "vault_kubernetes_auth_backend_role" "bookstore_mysql" {
  role_name                   = "bookstore-mysql-role"
  bound_service_account_names = var.mysql_service_accounts
  bound_service_account_namespaces = [var.namespace]
  token_policies = [vault_policy.bookstore_mysql.name]
  token_ttl                   = "3600" # 1 hour
}

resource "vault_kv_secret_v2" "minio" {
  for_each = toset(var.mysql_services)
  name  = "bookstore/${each.key}/mysql"
  mount = "kv"
  data_json = jsonencode({
    root_password = var.mysql_root_password,
    username      = var.mysql_username,
    password      = var.mysql_password
  })
}
