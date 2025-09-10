terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

resource "kubernetes_service_account_v1" "vso_oauth2_proxy" {
  metadata {
    name      = "vso-oauth2-proxy"
    namespace = var.namespace
  }
}

resource "vault_policy" "vso_oauth2_proxy" {
  name   = "vso-oauth2-proxy"
  policy = file("${path.module}/vso-policy.hcl")
}

resource "vault_kubernetes_auth_backend_role" "vso_oauth2_proxy" {
  role_name                        = "vso-oauth2-proxy-role"
  bound_service_account_names      = ["vso-oauth2-proxy"]
  bound_service_account_namespaces = [var.namespace]
  token_policies                   = [vault_policy.vso_oauth2_proxy.name]
  token_ttl                        = "3600"
}

resource "kubernetes_manifest" "vaultauth_oauth2_proxy" {
  manifest = {
    apiVersion = "secrets.hashicorp.com/v1beta1"
    kind       = "VaultAuth"
    metadata = {
      name      = "oauth2-proxy"
      namespace = var.namespace
    }
    spec = {
      method = "kubernetes"
      mount  = "kubernetes"
      kubernetes = {
        role           = "vso-oauth2-proxy-role"
        serviceAccount = kubernetes_service_account_v1.vso_oauth2_proxy.metadata[0].name
      }
    }
  }
}

resource "kubernetes_manifest" "vaultstaticsecret_oauth2_proxy" {
  manifest = {
    apiVersion = "secrets.hashicorp.com/v1beta1"
    kind       = "VaultStaticSecret"
    metadata = {
      name      = "oauth2-proxy-credentials"
      namespace = var.namespace
    }
    spec = {
      type         = "kv-v2"
      mount        = "kv"
      path         = var.vault_kv_path
      refreshAfter = "1h"
      vaultAuthRef = "oauth2-proxy"
      destination = {
        name   = var.secret_name
        create = true
      }
    }
  }
}
