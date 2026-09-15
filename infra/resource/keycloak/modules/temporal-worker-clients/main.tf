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

resource "keycloak_openid_client" "worker" {
  for_each    = toset(var.services)
  realm_id    = var.realm_id
  client_id   = "temporal-worker-${each.key}"
  name        = "Temporal Worker (${each.key})"
  description = "Confidential service-account client the ${each.key} service uses to authenticate its Temporal client and worker."
  enabled     = true

  access_type                  = "CONFIDENTIAL"
  client_authenticator_type    = "client-secret"
  standard_flow_enabled        = false
  implicit_flow_enabled        = false
  direct_access_grants_enabled = false
  service_accounts_enabled     = true
  full_scope_allowed           = false
}

resource "keycloak_openid_hardcoded_claim_protocol_mapper" "permissions" {
  for_each  = keycloak_openid_client.worker
  realm_id  = var.realm_id
  client_id = each.value.id
  name      = "temporal-permissions"

  claim_name          = var.permissions_claim_name
  claim_value         = jsonencode(var.permissions)
  claim_value_type    = "JSON"
  add_to_id_token     = false
  add_to_access_token = true
  add_to_userinfo     = false
}

resource "keycloak_openid_audience_protocol_mapper" "temporal" {
  for_each  = keycloak_openid_client.worker
  realm_id  = var.realm_id
  client_id = each.value.id
  name      = "temporal-audience"

  included_custom_audience = var.audience
  add_to_access_token      = true
  add_to_id_token          = false
}

resource "vault_kv_secret_v2" "worker" {
  for_each = keycloak_openid_client.worker
  mount    = "kv-bookstore"
  name     = "${each.key}/temporal"
  data_json = jsonencode({
    client_id     = each.value.client_id
    client_secret = each.value.client_secret
  })
}
