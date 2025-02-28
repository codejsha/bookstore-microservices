terraform {
  required_providers {
    gitea = {
      source = "go-gitea/gitea"
    }
  }
}

resource "gitea_team" "dev_team" {
  name                     = "developer-team"
  organisation             = var.org_name
  permission               = "write"
  include_all_repositories = false
  repositories             = var.dev_repos
}

resource "gitea_team_members" "dev_team_members" {
  team_id = gitea_team.dev_team.id
  members = var.dev_users
}

resource "gitea_token" "dev_admin" {
  name = "dev_admin"
  scopes = ["all"]
}

resource "gitea_team" "devops_team" {
  name                     = "devops-team"
  organisation             = var.org_name
  permission               = "write"
  include_all_repositories = false
  repositories             = var.devops_repos
}

resource "gitea_team_members" "devops_team_members" {
  team_id = gitea_team.devops_team.id
  members = var.devops_users
}
