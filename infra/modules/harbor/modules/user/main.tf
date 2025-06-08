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
    user.user_name => user
  }

  email     = "${each.value.user_name}@example.com"
  username  = each.value.user_name
  password  = each.value.user_password
  full_name = each.value.user_name
}
