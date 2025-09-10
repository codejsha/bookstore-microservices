terraform {
  required_providers {
    keycloak = {
      source = "keycloak/keycloak"
    }
  }
}

resource "keycloak_openid_client" "mobile" {
  realm_id    = var.realm_id
  client_id   = "bookstore-mobile"
  name        = "Bookstore Mobile (Native)"
  description = "Public client for the Expo / React Native app: authorization code flow with PKCE over the bookstore:// custom scheme."
  enabled     = true

  access_type                  = "PUBLIC"
  standard_flow_enabled        = true
  implicit_flow_enabled        = false
  direct_access_grants_enabled = false
  service_accounts_enabled     = false

  valid_redirect_uris             = var.valid_redirect_uris
  valid_post_logout_redirect_uris = var.valid_post_logout_redirect_uris

  pkce_code_challenge_method = "S256"
}
