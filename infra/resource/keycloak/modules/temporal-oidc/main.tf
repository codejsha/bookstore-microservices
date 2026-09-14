terraform {
  required_providers {
    keycloak = {
      source = "keycloak/keycloak"
    }
    vault = {
      source = "hashicorp/vault"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

resource "keycloak_openid_client" "temporal" {
  realm_id                     = var.realm_id
  client_id                    = "temporal"
  name                         = "Temporal"
  description                  = "Confidential client for Temporal Web UI SSO: OIDC login for the workflow console."
  enabled                      = true
  access_type                  = "CONFIDENTIAL"
  standard_flow_enabled        = true
  implicit_flow_enabled        = false
  direct_access_grants_enabled = false
  valid_redirect_uris = [
    "${var.temporal_url}/auth/sso/callback"
  ]
  web_origins = [var.temporal_url]
  root_url    = var.temporal_url
}

resource "keycloak_openid_client_default_scopes" "temporal" {
  realm_id       = var.realm_id
  client_id      = keycloak_openid_client.temporal.id
  default_scopes = concat(var.builtin_default_scopes, [var.groups_scope_name])
}

resource "vault_kv_secret_v2" "temporal_oidc_secret" {
  mount = "kv-infra"
  name  = "keycloak/temporal-oidc/client-secret"
  data_json = jsonencode({
    client_id     = keycloak_openid_client.temporal.client_id
    client_secret = keycloak_openid_client.temporal.client_secret
  })
}

resource "kubernetes_secret_v1" "temporal_oidc_secret" {
  metadata {
    name      = "temporal-web-oidc"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/part-of" = "temporal"
    }
  }
  data = {
    TEMPORAL_AUTH_CLIENT_ID     = keycloak_openid_client.temporal.client_id
    TEMPORAL_AUTH_CLIENT_SECRET = keycloak_openid_client.temporal.client_secret
  }
}

data "kubernetes_config_map_v1" "vault_pki_ca" {
  metadata {
    name      = var.ca_configmap_name
    namespace = var.ca_configmap_namespace
  }
}

resource "kubernetes_config_map_v1" "temporal_keycloak_ca" {
  metadata {
    name      = "temporal-keycloak-ca"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/part-of" = "temporal"
    }
  }
  data = {
    "ca.crt" = data.kubernetes_config_map_v1.vault_pki_ca.data[var.ca_configmap_key]
  }
}
