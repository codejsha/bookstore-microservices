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

locals {
  user_credentials_by_key = { for idx, user in var.user_credentials : "${user.username}-${idx}" => user }
}

resource "gitea_user" "users" {
  for_each   = nonsensitive(toset(keys(local.user_credentials_by_key)))
  username   = nonsensitive(local.user_credentials_by_key[each.key].username)
  login_name = nonsensitive(local.user_credentials_by_key[each.key].username)
  email      = "${nonsensitive(local.user_credentials_by_key[each.key].username)}@example.com"
  password   = local.user_credentials_by_key[each.key].password
}

resource "gitea_team" "team" {
  name                     = var.team_name
  organisation             = var.org_name
  permission               = "write"
  repositories             = var.user_repos
  include_all_repositories = false
}

resource "gitea_team_members" "team_members" {
  team_id = gitea_team.team.id
  members = [for user in gitea_user.users : user.username]
}

resource "vault_kv_secret_v2" "credentials" {
  for_each = nonsensitive(toset(keys(local.user_credentials_by_key)))
  name     = "gitea/users/${nonsensitive(local.user_credentials_by_key[each.key].username)}"
  mount    = "kv"
  data_json = jsonencode(
    {
      username = nonsensitive(local.user_credentials_by_key[each.key].username),
      password = local.user_credentials_by_key[each.key].password
    }
  )
}
