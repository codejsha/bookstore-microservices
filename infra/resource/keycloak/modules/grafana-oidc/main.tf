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

resource "keycloak_openid_client" "grafana" {
  realm_id                     = var.realm_id
  client_id                    = "grafana"
  name                         = "Grafana"
  description                  = "Confidential client for Grafana SSO: generic OAuth login with realm roles mapped to Grafana org roles."
  enabled                      = true
  access_type                  = "CONFIDENTIAL"
  standard_flow_enabled        = true
  implicit_flow_enabled        = false
  direct_access_grants_enabled = false
  valid_redirect_uris = [
    "${var.grafana_url}/login/generic_oauth"
  ]
  valid_post_logout_redirect_uris = [
    "${var.grafana_url}/login"
  ]
  web_origins = [var.grafana_url]
  root_url    = var.grafana_url

  pkce_code_challenge_method = "S256"
}

resource "keycloak_openid_user_realm_role_protocol_mapper" "realm_roles" {
  realm_id  = var.realm_id
  client_id = keycloak_openid_client.grafana.id
  name      = "realm-roles"

  claim_name          = "roles"
  multivalued         = true
  add_to_id_token     = true
  add_to_access_token = false
  add_to_userinfo     = true
}

resource "vault_kv_secret_v2" "grafana_oidc_secret" {
  mount = "kv-infra"
  name  = "keycloak/grafana-oidc/client-secret"
  data_json = jsonencode({
    client_secret = keycloak_openid_client.grafana.client_secret
  })
}

resource "kubernetes_secret_v1" "grafana_oidc_secret" {
  metadata {
    name      = "grafana-oidc"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/part-of" = "grafana"
    }
  }
  data = {
    client_secret = keycloak_openid_client.grafana.client_secret
  }
}

data "kubernetes_config_map_v1" "vault_pki_ca" {
  metadata {
    name      = var.ca_configmap_name
    namespace = var.ca_configmap_namespace
  }
}

resource "kubernetes_config_map_v1" "grafana_keycloak_ca" {
  metadata {
    name      = "grafana-keycloak-ca"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/part-of" = "grafana"
    }
  }
  data = {
    "ca.crt" = data.kubernetes_config_map_v1.vault_pki_ca.data[var.ca_configmap_key]
  }
}
