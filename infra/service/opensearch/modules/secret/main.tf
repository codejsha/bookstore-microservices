terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
    random = {
      source = "hashicorp/random"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

resource "random_password" "admin" {
  length           = 24
  special          = true
  override_special = "!#%^&*-_=+"
  min_upper        = 1
  min_lower        = 1
  min_numeric      = 1
  min_special      = 1
}

resource "vault_kv_secret_v2" "admin" {
  name  = "opensearch/admin/credentials"
  mount = "kv"
  data_json = jsonencode({
    username         = "admin"
    initial_password = random_password.admin.result
    password         = random_password.admin.result
  })
}

resource "vault_policy" "opensearch" {
  name   = "opensearch"
  policy = file("${path.module}/policy.hcl")
}

resource "vault_kubernetes_auth_backend_role" "opensearch" {
  role_name                        = "opensearch-role"
  bound_service_account_names      = ["opensearch"]
  bound_service_account_namespaces = [var.namespace]
  token_policies                   = [vault_policy.opensearch.name]
  token_ttl                        = "3600"
}

resource "kubernetes_service_account_v1" "vso_opensearch" {
  metadata {
    name      = "vso-opensearch"
    namespace = var.namespace
  }
}

resource "vault_policy" "vso_opensearch" {
  name   = "vso-opensearch"
  policy = file("${path.module}/vso-policy.hcl")
}

resource "vault_kubernetes_auth_backend_role" "vso_opensearch" {
  role_name                        = "vso-opensearch-role"
  bound_service_account_names      = ["vso-opensearch"]
  bound_service_account_namespaces = [var.namespace]
  token_policies                   = [vault_policy.vso_opensearch.name]
  token_ttl                        = "3600"
}

resource "kubernetes_manifest" "vaultauth_opensearch" {
  manifest = {
    apiVersion = "secrets.hashicorp.com/v1beta1"
    kind       = "VaultAuth"
    metadata = {
      name      = "opensearch"
      namespace = var.namespace
    }
    spec = {
      method = "kubernetes"
      mount  = "kubernetes"
      kubernetes = {
        role           = "vso-opensearch-role"
        serviceAccount = "vso-opensearch"
      }
    }
  }
}

resource "kubernetes_manifest" "vaultstaticsecret_opensearch_admin" {
  manifest = {
    apiVersion = "secrets.hashicorp.com/v1beta1"
    kind       = "VaultStaticSecret"
    metadata = {
      name      = "opensearch-admin"
      namespace = var.namespace
    }
    spec = {
      type         = "kv-v2"
      mount        = "kv"
      path         = "opensearch/admin/credentials"
      refreshAfter = "1h"
      vaultAuthRef = "opensearch"
      destination = {
        name   = "opensearch-admin-credentials"
        create = true
      }
    }
  }
}
