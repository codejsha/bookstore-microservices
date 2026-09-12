terraform {
  required_providers {
    keycloak = {
      source = "keycloak/keycloak"
    }
  }
}

resource "keycloak_role" "user" {
  realm_id    = var.realm_id
  name        = "USER"
  description = var.roles["USER"]
}

resource "keycloak_role" "staff" {
  realm_id        = var.realm_id
  name            = "STAFF"
  description     = var.roles["STAFF"]
  composite_roles = [keycloak_role.user.id]
}

resource "keycloak_role" "manage" {
  realm_id        = var.realm_id
  name            = "MANAGE"
  description     = var.roles["MANAGE"]
  composite_roles = [keycloak_role.staff.id]
}

resource "keycloak_role" "system" {
  realm_id    = var.realm_id
  name        = "SYSTEM"
  description = var.roles["SYSTEM"]
}

resource "keycloak_default_roles" "realm" {
  realm_id      = var.realm_id
  default_roles = concat(var.builtin_default_roles, [keycloak_role.user.name])
}

moved {
  from = keycloak_role.app["MANAGE"]
  to   = keycloak_role.manage
}

moved {
  from = keycloak_role.app["SYSTEM"]
  to   = keycloak_role.system
}
