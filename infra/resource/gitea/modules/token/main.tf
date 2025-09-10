terraform {
  required_providers {
    gitea = {
      source = "go-gitea/gitea"
    }
    vault = {
      source = "hashicorp/vault"
    }
  }
}

resource "gitea_token" "admin_token" {
  name   = "admin-token"
  scopes = ["all"]
}

resource "vault_kv_secret_v2" "admin_token" {
  name  = "gitea/admin/token"
  mount = "kv"
  data_json = jsonencode({
    token = gitea_token.admin_token.token
  })
}

resource "gitea_token" "resolver_token" {
  name   = "tekton-resolver-token"
  scopes = ["read:repository"]
}

resource "vault_kv_secret_v2" "resolver_token" {
  name  = "gitea/resolver/token"
  mount = "kv"
  data_json = jsonencode({
    token = gitea_token.resolver_token.token
  })
}
