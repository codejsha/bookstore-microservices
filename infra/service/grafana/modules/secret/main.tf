terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
    random = {
      source = "hashicorp/random"
    }
  }
}

resource "random_password" "admin" {
  length  = 32
  special = false
}

resource "vault_kv_secret_v2" "admin" {
  name  = "grafana/admin/credentials"
  mount = "kv"
  data_json = jsonencode({
    admin_user     = "admin"
    admin_password = random_password.admin.result
  })
}

resource "kubernetes_service_account_v1" "vso_grafana" {
  metadata {
    name      = "vso-grafana"
    namespace = var.namespace
  }
}

resource "vault_policy" "vso_grafana" {
  name   = "vso-grafana"
  policy = file("${path.module}/vso-policy.hcl")
}

resource "vault_kubernetes_auth_backend_role" "vso_grafana" {
  role_name                        = "vso-grafana-role"
  bound_service_account_names      = ["vso-grafana"]
  bound_service_account_namespaces = [var.namespace]
  token_policies                   = [vault_policy.vso_grafana.name]
  token_ttl                        = "3600"
}

resource "kubernetes_manifest" "vaultauth_grafana" {
  manifest = {
    apiVersion = "secrets.hashicorp.com/v1beta1"
    kind       = "VaultAuth"
    metadata = {
      name      = "grafana"
      namespace = var.namespace
    }
    spec = {
      method = "kubernetes"
      mount  = "kubernetes"
      kubernetes = {
        role           = "vso-grafana-role"
        serviceAccount = "vso-grafana"
      }
    }
  }
}

resource "kubernetes_manifest" "vaultstaticsecret_grafana_admin" {
  manifest = {
    apiVersion = "secrets.hashicorp.com/v1beta1"
    kind       = "VaultStaticSecret"
    metadata = {
      name      = "grafana-admin"
      namespace = var.namespace
    }
    spec = {
      type         = "kv-v2"
      mount        = "kv"
      path         = "grafana/admin/credentials"
      refreshAfter = "1h"
      vaultAuthRef = "grafana"
      destination = {
        name   = "grafana-admin-credentials"
        create = true
        annotations = {
          "reflector.v1.k8s.emberstack.com/reflection-allowed"            = "true"
          "reflector.v1.k8s.emberstack.com/reflection-allowed-namespaces" = "istio-system"
          "reflector.v1.k8s.emberstack.com/reflection-auto-enabled"       = "true"
          "reflector.v1.k8s.emberstack.com/reflection-auto-namespaces"    = "istio-system"
        }
      }
    }
  }
}
