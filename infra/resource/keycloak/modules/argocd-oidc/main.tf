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

resource "keycloak_openid_client" "argocd" {
  realm_id                     = var.realm_id
  client_id                    = "argocd"
  name                         = "ArgoCD"
  description                  = "Confidential client for Argo CD SSO: OIDC login with platform realm roles mapped into the groups claim."
  enabled                      = true
  access_type                  = "CONFIDENTIAL"
  standard_flow_enabled        = true
  implicit_flow_enabled        = false
  direct_access_grants_enabled = false
  valid_redirect_uris = [
    "${var.argocd_url}/auth/callback"
  ]
  web_origins = [var.argocd_url]
  root_url    = var.argocd_url
}

resource "keycloak_openid_client_default_scopes" "argocd" {
  realm_id       = var.realm_id
  client_id      = keycloak_openid_client.argocd.id
  default_scopes = concat(var.builtin_default_scopes, [var.groups_scope_name])
}

resource "vault_kv_secret_v2" "argocd_oidc_secret" {
  mount = "kv-infra"
  name  = "keycloak/argocd-oidc/client-secret"
  data_json = jsonencode({
    client_secret = keycloak_openid_client.argocd.client_secret
  })
}

resource "kubernetes_secret_v1" "argocd_oidc_secret" {
  metadata {
    name      = "argocd-oidc-secret"
    namespace = var.namespace
    labels = {
      "app.kubernetes.io/part-of" = "argocd"
    }
  }
  data = {
    "oidc.keycloak.clientSecret" = keycloak_openid_client.argocd.client_secret
  }
}
