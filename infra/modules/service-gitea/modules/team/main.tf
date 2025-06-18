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

resource "gitea_user" "users" {
  for_each   = {for idx, user in var.user_credentials : "${user.username}-${idx}" => user}
  username   = each.value.username
  login_name = each.value.username
  email      = "${each.value.username}@example.com"
  password   = each.value.password
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
  for_each = {for idx, user in var.user_credentials : "${user.username}-${idx}" => user}
  name     = "gitea/users/${each.value.username}"
  mount    = "kv"
  data_json = jsonencode(
    {
      username = each.value.username,
      password = each.value.password
    }
  )
}
