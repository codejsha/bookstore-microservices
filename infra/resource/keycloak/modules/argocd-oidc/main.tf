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
  description                  = "Confidential client for Argo CD SSO: OIDC login with realm group membership mapped into the groups claim."
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

resource "keycloak_group" "argocd_groups" {
  for_each = toset(["devadmin", "devopsadmin", "dev-user"])
  realm_id = var.realm_id
  name     = each.key
}

resource "keycloak_openid_group_membership_protocol_mapper" "argocd_groups" {
  realm_id   = var.realm_id
  client_id  = keycloak_openid_client.argocd.id
  name       = "groups"
  claim_name = "groups"
  full_path  = false
}

resource "vault_kv_secret_v2" "argocd_oidc_secret" {
  mount = "kv"
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
