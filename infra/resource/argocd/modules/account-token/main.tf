terraform {
  required_providers {
    argocd = {
      source = "argoproj-labs/argocd"
    }
    vault = {
      source = "hashicorp/vault"
    }
    time = {
      source = "hashicorp/time"
    }
  }
}

locals {
  accounts = ["dev-ci"]
}

resource "time_rotating" "account_token" {
  rotation_days = 90
}

resource "argocd_account_token" "account_tokens" {
  for_each   = toset(local.accounts)
  account    = each.key
  expires_in = "4320h"

  lifecycle {
    replace_triggered_by = [time_rotating.account_token]
  }
}

resource "vault_kv_secret_v2" "account_tokens" {
  for_each = toset(local.accounts)
  name     = "argocd/${each.key}/token"
  mount    = "kv"
  data_json = jsonencode({
    token = argocd_account_token.account_tokens[each.key].jwt
  })
}
