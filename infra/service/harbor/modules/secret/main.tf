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

resource "kubernetes_service_account_v1" "vso_harbor" {
  metadata {
    name      = "vso-harbor"
    namespace = var.namespace
  }
}

resource "vault_policy" "vso_harbor" {
  name   = "vso-harbor"
  policy = file("${path.module}/vso-policy.hcl")
}

resource "vault_kubernetes_auth_backend_role" "vso_harbor" {
  role_name                        = "vso-harbor-role"
  bound_service_account_names      = ["vso-harbor"]
  bound_service_account_namespaces = [var.namespace]
  token_policies                   = [vault_policy.vso_harbor.name]
  token_ttl                        = "3600"
}

resource "kubernetes_manifest" "vaultauth_harbor" {
  manifest = {
    apiVersion = "secrets.hashicorp.com/v1beta1"
    kind       = "VaultAuth"
    metadata = {
      name      = "harbor"
      namespace = var.namespace
    }
    spec = {
      method = "kubernetes"
      mount  = "kubernetes"
      kubernetes = {
        role           = "vso-harbor-role"
        serviceAccount = "vso-harbor"
      }
    }
  }
}

resource "kubernetes_manifest" "vaultstaticsecret_harbor_admin" {
  manifest = {
    apiVersion = "secrets.hashicorp.com/v1beta1"
    kind       = "VaultStaticSecret"
    metadata = {
      name      = "harbor-admin"
      namespace = var.namespace
    }
    spec = {
      type         = "kv-v2"
      mount        = "kv-infra"
      path         = "harbor/admin/credentials"
      refreshAfter = "1h"
      vaultAuthRef = "harbor"
      destination = {
        name   = "harbor-admin-credentials"
        create = true
      }
    }
  }
}

resource "kubernetes_manifest" "vaultstaticsecret_harbor_registry" {
  manifest = {
    apiVersion = "secrets.hashicorp.com/v1beta1"
    kind       = "VaultStaticSecret"
    metadata = {
      name      = "harbor-registry"
      namespace = var.namespace
    }
    spec = {
      type         = "kv-v2"
      mount        = "kv-infra"
      path         = "harbor/registry/credentials"
      refreshAfter = "1h"
      vaultAuthRef = "harbor"
      destination = {
        name   = "harbor-registry-credentials"
        create = true
      }
    }
  }
}

data "kubernetes_config_map_v1" "vault_pki_ca" {
  metadata {
    name      = var.ca_configmap_name
    namespace = var.ca_configmap_namespace
  }
}

resource "kubernetes_secret_v1" "ca_bundle" {
  metadata {
    name      = "harbor-ca-bundle"
    namespace = var.namespace
  }
  data = {
    "ca.crt" = data.kubernetes_config_map_v1.vault_pki_ca.data[var.ca_configmap_key]
  }
}
