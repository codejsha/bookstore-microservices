terraform {
  required_providers {
    harbor = {
      source = "goharbor/harbor"
    }
  }
}

resource "harbor_user" "dev_admin" {
  email     = "devadmin@example.com"
  username  = "devadmin"
  password  = var.admin_password
  full_name = "devadmin"
}

resource "harbor_user" "devops_admin" {
  email     = "devopsadmin@example.com"
  username  = "devopsadmin"
  password  = var.admin_password
  full_name = "devopsadmin"
}
