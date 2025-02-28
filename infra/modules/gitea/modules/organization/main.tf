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

resource "gitea_org" "organization" {
  name = var.org_name
}
