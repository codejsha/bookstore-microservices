terraform {
  required_providers {
    keycloak = {
      source = "keycloak/keycloak"
    }
  }
}

resource "keycloak_openid_client" "web" {
  realm_id    = var.realm_id
  client_id   = "bookstore-web"
  name        = "Bookstore Web (SPA)"
  description = "Public client for the customer storefront SPA: authorization code flow with PKCE from the browser."
  enabled     = true

  access_type                  = "PUBLIC"
  standard_flow_enabled        = true
  implicit_flow_enabled        = false
  direct_access_grants_enabled = false
  service_accounts_enabled     = false

  root_url                        = var.root_url
  valid_redirect_uris             = var.valid_redirect_uris
  valid_post_logout_redirect_uris = var.valid_post_logout_redirect_uris
  web_origins                     = var.web_origins

  pkce_code_challenge_method = "S256"
}
