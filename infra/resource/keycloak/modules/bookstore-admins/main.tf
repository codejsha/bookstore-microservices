terraform {
  required_providers {
    keycloak = {
      source = "keycloak/keycloak"
    }
    random = {
      source = "hashicorp/random"
    }
  }
}

resource "random_password" "account" {
  for_each = var.accounts
  length   = 32
  special  = false
}

resource "keycloak_user" "account" {
  for_each       = var.accounts
  realm_id       = var.realm_id
  username       = each.value.username
  email          = each.value.email
  email_verified = true
  first_name     = each.value.first_name
  last_name      = each.value.last_name
  enabled        = true

  initial_password {
    value     = random_password.account[each.key].result
    temporary = true
  }

  lifecycle {
    ignore_changes = [required_actions]
  }
}

resource "keycloak_user_roles" "account" {
  for_each = var.accounts
  realm_id = var.realm_id
  user_id  = keycloak_user.account[each.key].id
  role_ids = [var.role_ids[each.value.role]]
}
