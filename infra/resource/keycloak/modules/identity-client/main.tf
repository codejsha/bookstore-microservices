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

resource "random_password" "identity_admin" {
  length  = 32
  special = false
}

resource "keycloak_user" "identity_admin" {
  realm_id       = var.realm_id
  username       = var.identity_admin_username
  email          = var.identity_admin_email
  email_verified = true
  first_name     = "Identity"
  last_name      = "Service"
  enabled        = true

  initial_password {
    value     = random_password.identity_admin.result
    temporary = false
  }
}

resource "keycloak_user_roles" "identity_admin" {
  realm_id = var.realm_id
  user_id  = keycloak_user.identity_admin.id
  role_ids = [data.keycloak_role.realm_admin.id]
}

resource "random_password" "manager" {
  length  = 32
  special = false
}

resource "keycloak_user" "manager" {
  realm_id       = var.realm_id
  username       = var.manager_username
  email          = var.manager_username
  email_verified = true
  first_name     = var.manager_first_name
  last_name      = var.manager_last_name
  enabled        = true

  initial_password {
    value     = random_password.manager.result
    temporary = false
  }
}

resource "keycloak_user_roles" "manager" {
  realm_id = var.realm_id
  user_id  = keycloak_user.manager.id
  role_ids = [var.manage_role_id]
}

resource "vault_kv_secret_v2" "manager" {
  mount = "kv-infra"
  name  = "keycloak/bookstore/manager/credentials"
  data_json = jsonencode({
    username = keycloak_user.manager.username
    password = random_password.manager.result
  })
}

moved {
  from = random_password.realm_admin
  to   = random_password.manager
}

moved {
  from = keycloak_user.realm_admin
  to   = keycloak_user.manager
}

moved {
  from = keycloak_user_roles.realm_admin_roles
  to   = keycloak_user_roles.manager
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
  mount = "kv-bookstore"
  name  = "identity/keycloak"
  data_json = jsonencode({
    client_id      = keycloak_openid_client.identity.client_id
    client_secret  = keycloak_openid_client.identity.client_secret
    admin_username = keycloak_user.identity_admin.username
    admin_password = random_password.identity_admin.result
  })
}
