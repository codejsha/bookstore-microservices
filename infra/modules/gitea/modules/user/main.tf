terraform {
  required_providers {
    gitea = {
      source = "go-gitea/gitea"
    }
  }
}

resource "gitea_user" "dev_admin" {
  username   = "devadmin"
  login_name = "devadmin"
  email      = "devadmin@example.com"
  password   = var.admin_password
}

resource "gitea_user" "devops_admin" {
  username   = "devopsadmin"
  login_name = "devopsadmin"
  email      = "devopsadmin@example.com"
  password   = var.admin_password
}
