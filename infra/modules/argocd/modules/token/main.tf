terraform {
  required_providers {
    argocd = {
      source = "argoproj-labs/argocd"
    }
    vault = {
      source = "hashicorp/vault"
    }
  }
}

locals {
  accounts = ["devadmin", "devopsadmin"]
}

resource "argocd_account_token" "account_tokens" {
  for_each = toset(local.accounts)
  account    = each.key
  expires_in = "8760h"
}

resource "vault_kv_secret_v2" "account_tokens" {
  for_each = toset(local.accounts)
  name  = "argocd/${each.key}/token"
  mount = "kv"
  data_json = jsonencode({
    token = argocd_account_token.account_tokens[each.key].jwt
  })
}
