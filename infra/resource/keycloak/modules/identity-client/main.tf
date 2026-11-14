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

resource "keycloak_openid_client" "identity" {
  realm_id    = var.realm_id
  client_id   = "identity"
  name        = "Identity Service"
  description = "Confidential client for the identity service: signin/signup/token endpoints and Keycloak user administration."
  enabled     = true

  access_type                  = "CONFIDENTIAL"
  client_authenticator_type    = "client-secret"
  standard_flow_enabled        = false
  implicit_flow_enabled        = false
  direct_access_grants_enabled = true
  service_accounts_enabled     = true
}

resource "keycloak_openid_client_service_account_role" "identity" {
  for_each                = toset(var.service_account_roles)
  realm_id                = var.realm_id
  service_account_user_id = keycloak_openid_client.identity.service_account_user_id
  client_id               = data.keycloak_openid_client.realm_management.id
  role                    = each.value
}

data "keycloak_openid_client" "realm_management" {
  realm_id  = var.realm_id
  client_id = "realm-management"
}

resource "vault_kv_secret_v2" "keycloak_admin" {
  mount = "kv-bookstore"
  name  = "keycloak/admin/credentials"
  data_json = jsonencode({
    client_id     = keycloak_openid_client.identity.client_id
    client_secret = keycloak_openid_client.identity.client_secret
  })
}
