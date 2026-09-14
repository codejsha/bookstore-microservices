terraform {
  required_providers {
    gitea = {
      source = "go-gitea/gitea"
    }
  }
}

resource "gitea_team" "team" {
  name                     = var.team_name
  organisation             = var.org_name
  permission               = "write"
  repositories             = var.user_repos
  include_all_repositories = false
}

