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

resource "kubernetes_service_account_v1" "vso_hyperswitch" {
  metadata {
    name      = "vso-hyperswitch"
    namespace = var.namespace
  }
}

resource "vault_policy" "vso_hyperswitch" {
  name   = "vso-hyperswitch"
  policy = file("${path.module}/vso-policy.hcl")
}

resource "vault_kubernetes_auth_backend_role" "vso_hyperswitch" {
  role_name                        = "vso-hyperswitch-role"
  bound_service_account_names      = ["vso-hyperswitch"]
  bound_service_account_namespaces = [var.namespace]
  token_policies                   = [vault_policy.vso_hyperswitch.name]
  token_ttl                        = "3600"
}

resource "kubernetes_manifest" "vaultauth_hyperswitch" {
  manifest = {
    apiVersion = "secrets.hashicorp.com/v1beta1"
    kind       = "VaultAuth"
    metadata = {
      name      = "hyperswitch"
      namespace = var.namespace
    }
    spec = {
      method = "kubernetes"
      mount  = "kubernetes"
      kubernetes = {
        role           = "vso-hyperswitch-role"
        serviceAccount = "vso-hyperswitch"
      }
    }
  }
}

resource "kubernetes_manifest" "vaultstaticsecret_hyperswitch_app" {
  manifest = {
    apiVersion = "secrets.hashicorp.com/v1beta1"
    kind       = "VaultStaticSecret"
    metadata = {
      name      = "hyperswitch-app"
      namespace = var.namespace
    }
    spec = {
      type         = "kv-v2"
      mount        = "kv"
      path         = "hyperswitch/app/credentials"
      refreshAfter = "1h"
      vaultAuthRef = "hyperswitch"
      destination = {
        name   = "hyperswitch-app-credentials"
        create = true
      }
    }
  }
}
