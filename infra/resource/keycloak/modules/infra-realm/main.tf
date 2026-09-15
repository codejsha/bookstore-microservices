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

resource "keycloak_realm" "infra" {
  realm        = var.realm_name
  enabled      = true
  display_name = "Platform Infra"

  login_theme   = "keycloak"
  account_theme = "keycloak.v3"

  registration_allowed = false
  ssl_required         = "external"

  access_token_lifespan    = "5m"
  sso_session_idle_timeout = "30m"
  sso_session_max_lifespan = "8h"

  security_defenses {
    brute_force_detection {
      permanent_lockout                = false
      max_login_failures               = 3
      wait_increment_seconds           = 300
      quick_login_check_milli_seconds  = 1000
      minimum_quick_login_wait_seconds = 300
      max_failure_wait_seconds         = 3600
      failure_reset_time_seconds       = 43200
    }
  }
}

resource "keycloak_required_action" "configure_totp" {
  realm_id       = keycloak_realm.infra.id
  alias          = "CONFIGURE_TOTP"
  name           = "Configure OTP"
  enabled        = true
  default_action = true
}

resource "keycloak_role" "developer" {
  realm_id        = keycloak_realm.infra.id
  name            = "DEVELOPER"
  description     = var.roles["DEVELOPER"]
  composite_roles = lookup(var.role_composites, "DEVELOPER", [])
}

resource "keycloak_role" "manager" {
  realm_id        = keycloak_realm.infra.id
  name            = "MANAGER"
  description     = var.roles["MANAGER"]
  composite_roles = lookup(var.role_composites, "MANAGER", [])
}

resource "keycloak_role" "admin" {
  realm_id        = keycloak_realm.infra.id
  name            = "ADMIN"
  description     = var.roles["ADMIN"]
  composite_roles = lookup(var.role_composites, "ADMIN", [])
}

moved {
  from = keycloak_role.operator
  to   = keycloak_role.manager
}

locals {
  role_ids = {
    DEVELOPER = keycloak_role.developer.id
    MANAGER   = keycloak_role.manager.id
    ADMIN     = keycloak_role.admin.id
  }
}

resource "random_password" "bootstrap" {
  for_each = var.bootstrap_accounts
  length   = 32
  special  = false
}

resource "keycloak_user" "bootstrap" {
  for_each       = var.bootstrap_accounts
  realm_id       = keycloak_realm.infra.id
  username       = each.value.username
  email          = each.value.email
  first_name     = each.value.first_name
  last_name      = each.value.last_name
  email_verified = true
  enabled        = true

  required_actions = [keycloak_required_action.configure_totp.alias]

  initial_password {
    value     = random_password.bootstrap[each.key].result
    temporary = false
  }

  lifecycle {
    ignore_changes = [required_actions]
  }
}

resource "keycloak_user_roles" "bootstrap" {
  for_each   = var.bootstrap_accounts
  realm_id   = keycloak_realm.infra.id
  user_id    = keycloak_user.bootstrap[each.key].id
  role_ids   = [local.role_ids[each.value.role]]
  exhaustive = false
}

resource "vault_kv_secret_v2" "bootstrap" {
  for_each = var.bootstrap_accounts
  mount    = "kv-bookstore"
  name     = "admin/${each.value.username}/credentials"
  data_json = jsonencode({
    username = keycloak_user.bootstrap[each.key].username
    password = random_password.bootstrap[each.key].result
  })
}

moved {
  from = random_password.operator
  to   = random_password.bootstrap["operator"]
}

moved {
  from = random_password.developer
  to   = random_password.bootstrap["developer"]
}

moved {
  from = keycloak_user.operator
  to   = keycloak_user.bootstrap["operator"]
}

moved {
  from = keycloak_user.developer
  to   = keycloak_user.bootstrap["developer"]
}

moved {
  from = keycloak_user_roles.operator
  to   = keycloak_user_roles.bootstrap["operator"]
}

moved {
  from = keycloak_user_roles.developer
  to   = keycloak_user_roles.bootstrap["developer"]
}

resource "keycloak_openid_client_scope" "groups" {
  realm_id               = keycloak_realm.infra.id
  name                   = "groups"
  description            = "Platform realm roles exposed as the groups claim for tools that authorize by group"
  include_in_token_scope = true
}

resource "keycloak_openid_user_realm_role_protocol_mapper" "groups" {
  realm_id        = keycloak_realm.infra.id
  client_scope_id = keycloak_openid_client_scope.groups.id
  name            = "realm-roles-as-groups"

  claim_name          = "groups"
  multivalued         = true
  add_to_id_token     = true
  add_to_access_token = false
  add_to_userinfo     = true
}
