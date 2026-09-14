terraform {
  required_providers {
    keycloak = {
      source = "keycloak/keycloak"
    }
    vault = {
      source = "hashicorp/vault"
    }
  }
}

resource "keycloak_openid_client" "harbor" {
  realm_id                     = var.realm_id
  client_id                    = "harbor"
  name                         = "Harbor"
  description                  = "Confidential client for Harbor SSO: OIDC login with platform realm roles mapped into the groups claim."
  enabled                      = true
  access_type                  = "CONFIDENTIAL"
  standard_flow_enabled        = true
  implicit_flow_enabled        = false
  direct_access_grants_enabled = false
  valid_redirect_uris = [
    "${var.harbor_url}/c/oidc/callback"
  ]
  web_origins = [var.harbor_url]
  root_url    = var.harbor_url
}

resource "keycloak_openid_client_default_scopes" "harbor" {
  realm_id       = var.realm_id
  client_id      = keycloak_openid_client.harbor.id
  default_scopes = concat(var.builtin_default_scopes, [var.groups_scope_name])
}

resource "vault_kv_secret_v2" "harbor_oidc_secret" {
  mount = "kv-infra"
  name  = "keycloak/harbor-oidc/client-secret"
  data_json = jsonencode({
    client_id     = keycloak_openid_client.harbor.client_id
    client_secret = keycloak_openid_client.harbor.client_secret
  })
}
