terraform {
  required_providers {
    harbor = {
      source = "goharbor/harbor"
    }
  }
}

locals {
  harbor_users_by_name = { for u in var.harbor_users : u.username => u }
}

resource "harbor_user" "harbor_users" {
  for_each            = nonsensitive(toset([for u in var.harbor_users : u.username]))
  username            = each.value
  email               = nonsensitive(local.harbor_users_by_name[each.value].email)
  password_wo         = local.harbor_users_by_name[each.value].password
  password_wo_version = 1
  full_name           = each.value
}
