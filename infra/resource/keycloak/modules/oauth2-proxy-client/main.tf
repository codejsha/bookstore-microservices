terraform {
  required_providers {
    keycloak = {
      source = "keycloak/keycloak"
    }
    vault = {
      source = "hashicorp/vault"
    }
    random = {
      source = "hashicorp/random"
    }
  }
}

resource "keycloak_openid_client" "oauth2_proxy" {
  realm_id    = var.realm_id
  client_id   = var.client_id
  name        = "Bookstore Edge (oauth2-proxy BFF)"
  description = "Confidential client for the oauth2-proxy edge BFF: cookie-session login for the admin console and storefront edge routes."
  enabled     = true

  access_type                  = "CONFIDENTIAL"
  client_authenticator_type    = "client-secret"
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

resource "keycloak_openid_user_realm_role_protocol_mapper" "realm_roles_in_id_token" {
  realm_id  = var.realm_id
  client_id = keycloak_openid_client.oauth2_proxy.id
  name      = "realm-roles-in-id-token"

  claim_name          = "realm_access.roles"
  multivalued         = true
  add_to_id_token     = true
  add_to_access_token = false
  add_to_userinfo     = false
}

resource "random_password" "cookie_secret" {
  length  = 32
  special = false
}

resource "vault_kv_secret_v2" "oauth2_proxy" {
  mount = "kv"
  name  = "bookstore/oauth2-proxy/config"
  data_json = jsonencode({
    "client-id"     = keycloak_openid_client.oauth2_proxy.client_id
    "client-secret" = keycloak_openid_client.oauth2_proxy.client_secret
    "cookie-secret" = random_password.cookie_secret.result
  })
}
