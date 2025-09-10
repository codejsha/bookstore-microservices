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

resource "vault_policy" "nexus" {
  name   = "nexus"
  policy = file("${path.module}/policy.hcl")
}

resource "vault_kubernetes_auth_backend_role" "nexus" {
  role_name                        = "nexus-role"
  bound_service_account_names      = ["nexus"]
  bound_service_account_namespaces = [var.namespace]
  token_policies                   = [vault_policy.nexus.name]
  token_ttl                        = "3600"
}

resource "kubernetes_manifest" "vaultauth_nexus" {
  manifest = {
    apiVersion = "secrets.hashicorp.com/v1beta1"
    kind       = "VaultAuth"
    metadata = {
      name      = "nexus"
      namespace = var.namespace
    }
    spec = {
      method = "kubernetes"
      mount  = "kubernetes"
      kubernetes = {
        role           = "nexus-role"
        serviceAccount = "nexus"
      }
    }
  }
}

resource "kubernetes_manifest" "vaultstaticsecret_nexus_admin" {
  manifest = {
    apiVersion = "secrets.hashicorp.com/v1beta1"
    kind       = "VaultStaticSecret"
    metadata = {
      name      = "nexus-admin-credentials"
      namespace = var.namespace
    }
    spec = {
      type         = "kv-v2"
      mount        = "kv"
      path         = "nexus/admin/credentials"
      refreshAfter = "1h"
      vaultAuthRef = "nexus"
      destination = {
        name   = "nexus-admin-credentials"
        create = true
      }
    }
  }
}
