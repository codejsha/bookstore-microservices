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

resource "random_password" "realm_admin" {
  length  = 32
  special = false
}

resource "keycloak_user" "realm_admin" {
  realm_id       = var.realm_id
  username       = var.realm_admin_username
  email          = var.realm_admin_username
  email_verified = true
  first_name     = var.realm_admin_first_name
  last_name      = var.realm_admin_last_name
  enabled        = true

  initial_password {
    value     = random_password.realm_admin.result
    temporary = false
  }
}

resource "keycloak_user_roles" "realm_admin_roles" {
  realm_id = var.realm_id
  user_id  = keycloak_user.realm_admin.id
  role_ids = [
    data.keycloak_role.realm_admin.id,
    var.manage_role_id,
  ]
}

data "keycloak_role" "realm_admin" {
  realm_id  = var.realm_id
  name      = "realm-admin"
  client_id = data.keycloak_openid_client.realm_management.id
}

data "keycloak_openid_client" "realm_management" {
  realm_id  = var.realm_id
  client_id = "realm-management"
}

resource "vault_kv_secret_v2" "identity_keycloak" {
  mount = "kv"
  name  = "bookstore/identity/keycloak"
  data_json = jsonencode({
    client_id      = keycloak_openid_client.identity.client_id
    client_secret  = keycloak_openid_client.identity.client_secret
    admin_username = var.realm_admin_username
    admin_password = random_password.realm_admin.result
  })
}
