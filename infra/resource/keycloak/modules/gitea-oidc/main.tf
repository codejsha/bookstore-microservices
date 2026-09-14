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

resource "keycloak_openid_client" "gitea" {
  realm_id                     = var.realm_id
  client_id                    = "gitea"
  name                         = "Gitea"
  description                  = "Confidential client for Gitea SSO: OIDC login with platform realm roles mapped into the groups claim."
  enabled                      = true
  access_type                  = "CONFIDENTIAL"
  standard_flow_enabled        = true
  implicit_flow_enabled        = false
  direct_access_grants_enabled = false
  valid_redirect_uris = [
    "${var.gitea_url}/user/oauth2/${var.auth_source_name}/callback"
  ]
  web_origins = [var.gitea_url]
  root_url    = var.gitea_url
}

resource "keycloak_openid_client_default_scopes" "gitea" {
  realm_id       = var.realm_id
  client_id      = keycloak_openid_client.gitea.id
  default_scopes = concat(var.builtin_default_scopes, [var.groups_scope_name])
}

resource "vault_kv_secret_v2" "gitea_oidc_secret" {
  mount = "kv-infra"
  name  = "keycloak/gitea-oidc/client-secret"
  data_json = jsonencode({
    client_id     = keycloak_openid_client.gitea.client_id
    client_secret = keycloak_openid_client.gitea.client_secret
  })
}

resource "kubernetes_secret_v1" "gitea_oidc_secret" {
  metadata {
    name      = "gitea-oidc"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/part-of" = "gitea"
    }
  }
  data = {
    key    = keycloak_openid_client.gitea.client_id
    secret = keycloak_openid_client.gitea.client_secret
  }
}

data "kubernetes_config_map_v1" "vault_pki_ca" {
  metadata {
    name      = var.ca_configmap_name
    namespace = var.ca_configmap_namespace
  }
}

resource "kubernetes_config_map_v1" "gitea_keycloak_ca" {
  metadata {
    name      = "gitea-keycloak-ca"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/part-of" = "gitea"
    }
  }
  data = {
    "ca.crt" = data.kubernetes_config_map_v1.vault_pki_ca.data[var.ca_configmap_key]
  }
}
