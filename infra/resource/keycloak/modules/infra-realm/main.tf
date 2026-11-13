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
  enabled        = false
  default_action = false
}

resource "keycloak_role" "role" {
  for_each        = var.roles
  realm_id        = keycloak_realm.infra.id
  name            = each.key
  description     = each.value
  composite_roles = lookup(var.role_composites, each.key, [])
}

moved {
  from = keycloak_role.operator
  to   = keycloak_role.manager
}

moved {
  from = keycloak_role.developer
  to   = keycloak_role.role["DEVELOPER"]
}

moved {
  from = keycloak_role.manager
  to   = keycloak_role.role["MANAGER"]
}

moved {
  from = keycloak_role.admin
  to   = keycloak_role.role["ADMIN"]
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
  role_ids   = [for role in each.value.roles : keycloak_role.role[role].id]
  exhaustive = false
}

resource "vault_kv_secret_v2" "credentials" {
  for_each = { for key, account in var.bootstrap_accounts : key => account if account.store_in_vault }
  mount    = "kv-infra"
  name     = "keycloak/admin/${each.value.username}/credentials"
  data_json = jsonencode({
    username = keycloak_user.bootstrap[each.key].username
    password = random_password.bootstrap[each.key].result
  })
}

removed {
  from = vault_kv_secret_v2.bootstrap

  lifecycle {
    destroy = false
  }
}

moved {
  from = random_password.bootstrap["platform_admin"]
  to   = random_password.bootstrap["bookstore_admin"]
}

moved {
  from = keycloak_user.bootstrap["platform_admin"]
  to   = keycloak_user.bootstrap["bookstore_admin"]
}

moved {
  from = keycloak_user_roles.bootstrap["platform_admin"]
  to   = keycloak_user_roles.bootstrap["bookstore_admin"]
}

moved {
  from = vault_kv_secret_v2.credentials["platform_admin"]
  to   = vault_kv_secret_v2.credentials["bookstore_admin"]
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
