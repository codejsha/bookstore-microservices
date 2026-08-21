terraform {
  required_providers {
    harbor = {
      source = "goharbor/harbor"
    }
  }
}

resource "harbor_user" "harbor_users" {
  for_each = {
    for user in var.harbor_users :
    user.username => user
  }
  email     = "${each.value.username}@example.com"
  username  = each.value.username
  password  = each.value.password
  full_name = each.value.username
}
