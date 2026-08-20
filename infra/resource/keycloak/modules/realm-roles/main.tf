terraform {
  required_providers {
    keycloak = {
      source = "keycloak/keycloak"
    }
  }
}

resource "keycloak_role" "app" {
  for_each = var.roles

  realm_id    = var.realm_id
  name        = each.key
  description = each.value
}
