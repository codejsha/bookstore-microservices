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
    random = {
      source = "hashicorp/random"
    }
  }
}

resource "keycloak_openid_client" "opensearch" {
  realm_id                     = var.realm_id
  client_id                    = "opensearch-dashboards"
  name                         = "OpenSearch Dashboards"
  description                  = "Confidential client for OpenSearch Dashboards SSO: OIDC login with platform realm roles mapped into the groups claim as security plugin backend roles."
  enabled                      = true
  access_type                  = "CONFIDENTIAL"
  standard_flow_enabled        = true
  implicit_flow_enabled        = false
  direct_access_grants_enabled = false
  valid_redirect_uris = [
    "${var.opensearch_url}/auth/openid/login"
  ]
  valid_post_logout_redirect_uris = [
    var.opensearch_url,
    "${var.opensearch_url}/*"
  ]
  web_origins = [var.opensearch_url]
  root_url    = var.opensearch_url
}

resource "keycloak_openid_client_default_scopes" "opensearch" {
  realm_id       = var.realm_id
  client_id      = keycloak_openid_client.opensearch.id
  default_scopes = concat(var.builtin_default_scopes, [var.groups_scope_name])
}

resource "random_password" "cookie" {
  length  = 48
  special = false
}

resource "vault_kv_secret_v2" "opensearch_oidc_secret" {
  mount = "kv-infra"
  name  = "keycloak/opensearch-oidc/client-secret"
  data_json = jsonencode({
    client_id       = keycloak_openid_client.opensearch.client_id
    client_secret   = keycloak_openid_client.opensearch.client_secret
    cookie_password = random_password.cookie.result
  })
}

resource "kubernetes_secret_v1" "opensearch_oidc_secret" {
  metadata {
    name      = "opensearch-dashboards-oidc"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/part-of" = "opensearch"
    }
  }
  data = {
    OPENSEARCH_OIDC_CLIENT_ID     = keycloak_openid_client.opensearch.client_id
    OPENSEARCH_OIDC_CLIENT_SECRET = keycloak_openid_client.opensearch.client_secret
    OPENSEARCH_COOKIE_PASSWORD    = random_password.cookie.result
  }
}

data "kubernetes_config_map_v1" "vault_pki_ca" {
  metadata {
    name      = var.ca_configmap_name
    namespace = var.ca_configmap_namespace
  }
}

resource "kubernetes_config_map_v1" "opensearch_keycloak_ca" {
  metadata {
    name      = "opensearch-keycloak-ca"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/part-of" = "opensearch"
    }
  }
  data = {
    "ca.crt" = data.kubernetes_config_map_v1.vault_pki_ca.data[var.ca_configmap_key]
  }
}
