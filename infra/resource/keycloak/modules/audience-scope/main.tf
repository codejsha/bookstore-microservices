terraform {
  required_providers {
    keycloak = {
      source = "keycloak/keycloak"
    }
  }
}

resource "keycloak_openid_client_scope" "bookstore_audience" {
  realm_id    = var.realm_id
  name        = "bookstore-audience"
  description = "Stamps a fixed aud claim consumed by the service mesh RequestAuthentication"

  include_in_token_scope = false
}

resource "keycloak_openid_audience_protocol_mapper" "bookstore_audience" {
  realm_id        = var.realm_id
  client_scope_id = keycloak_openid_client_scope.bookstore_audience.id
  name            = "bookstore-audience"

  included_custom_audience = var.audience
  add_to_access_token      = true
  add_to_id_token          = true
}

resource "keycloak_openid_user_realm_role_protocol_mapper" "bookstore_roles" {
  realm_id        = var.realm_id
  client_scope_id = keycloak_openid_client_scope.bookstore_audience.id
  name            = "bookstore-realm-roles"

  claim_name          = "roles"
  multivalued         = true
  add_to_access_token = true
  add_to_id_token     = true
  add_to_userinfo     = false
}

locals {
  base_default_scopes = ["acr", "basic", "email", "profile", "roles", "web-origins"]
}

resource "keycloak_openid_client_default_scopes" "bookstore_audience" {
  for_each = var.clients

  realm_id  = var.realm_id
  client_id = each.value

  default_scopes = concat(
    local.base_default_scopes,
    [keycloak_openid_client_scope.bookstore_audience.name],
    contains(var.service_account_clients, each.key) ? ["service_account"] : [],
  )
}
