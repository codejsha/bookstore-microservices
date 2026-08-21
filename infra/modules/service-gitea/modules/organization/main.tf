terraform {
  required_providers {
    gitea = {
      source = "go-gitea/gitea"
    }
  }
}

resource "gitea_org" "organization" {
  name = var.org_name
}
