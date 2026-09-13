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

resource "keycloak_openid_client" "kiali" {
  realm_id                     = var.realm_id
  client_id                    = "kiali"
  name                         = "Kiali"
  description                  = "Confidential client for Kiali SSO: OIDC login for the service mesh console."
  enabled                      = true
  access_type                  = "CONFIDENTIAL"
  standard_flow_enabled        = true
  implicit_flow_enabled        = false
  direct_access_grants_enabled = false
  valid_redirect_uris = [
    "${var.kiali_url}/kiali/*"
  ]
  web_origins = [var.kiali_url]
  root_url    = var.kiali_url
}

resource "vault_kv_secret_v2" "kiali_oidc_secret" {
  mount = "kv"
  name  = "keycloak/kiali-oidc/client-secret"
  data_json = jsonencode({
    client_secret = keycloak_openid_client.kiali.client_secret
  })
}

resource "kubernetes_secret_v1" "kiali_oidc_secret" {
  metadata {
    name      = "kiali"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/part-of" = "kiali"
    }
  }
  data = {
    "oidc-secret" = keycloak_openid_client.kiali.client_secret
  }
}

data "kubernetes_config_map_v1" "vault_pki_ca" {
  metadata {
    name      = var.ca_configmap_name
    namespace = var.ca_configmap_namespace
  }
}

resource "kubernetes_config_map_v1" "kiali_cabundle" {
  metadata {
    name      = "kiali-cabundle"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/part-of" = "kiali"
    }
  }
  data = {
    "openid-server-ca.crt" = data.kubernetes_config_map_v1.vault_pki_ca.data[var.ca_configmap_key]
  }
}
