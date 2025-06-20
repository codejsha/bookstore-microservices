terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
  }
}

resource "vault_pki_secret_backend_role" "opentelemetry" {
  name    = "opentelemetry"
  backend = "pki_int"
  allowed_domains = [
    "localhost",
    "opentelemetry-opentelemetry-operator-webhook",
    "opentelemetry-opentelemetry-operator-webhook.opentelemetry.svc",
    "opentelemetry-opentelemetry-operator-webhook.opentelemetry.svc.cluster.local"
  ]
  allow_bare_domains = true
  allow_glob_domains = true
  require_cn         = false
  max_ttl            = "1440h"
}

resource "kubernetes_service_account" "opentelemetry_issuer" {
  metadata {
    name      = "opentelemetry-issuer"
    namespace = var.namespace
  }
}

resource "kubernetes_secret" "opentelemetry_issuer_token" {
  type = "kubernetes.io/service-account-token"
  metadata {
    name      = "opentelemetry-issuer-token"
    namespace = var.namespace
    annotations = {
      "kubernetes.io/service-account.name" = kubernetes_service_account.opentelemetry_issuer.metadata[0].name
    }
  }
}

resource "vault_kubernetes_auth_backend_role" "opentelemetry_issuer" {
  role_name = "opentelemetry-issuer"
  backend   = "kubernetes"
  bound_service_account_names = [kubernetes_service_account.opentelemetry_issuer.metadata[0].name]
  bound_service_account_namespaces = [var.namespace]
  token_policies = ["pki_int"]
  token_ttl = "3600" # 1 hour
}

resource "kubernetes_cluster_role_binding" "opentelemetry_issuer_token_rolebinding" {
  metadata {
    name = "opentelemetry-issuer-token-rolebinding"
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = "system:auth-delegator"
  }
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.opentelemetry_issuer.metadata[0].name
    namespace = var.namespace
  }
}

resource "kubernetes_manifest" "opentelemetry_issuer" {
  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "Issuer"
    metadata = {
      name      = "opentelemetry-issuer"
      namespace = var.namespace
    }
    spec = {
      vault = {
        server = "http://vault.vault.svc.cluster.local:8200"
        path   = "pki_int/sign/${vault_pki_secret_backend_role.opentelemetry.name}"
        caBundle = base64encode(trimspace(file(var.kube_ca_crt_path)))
        auth = {
          kubernetes = {
            mountPath = "/v1/auth/kubernetes"
            role      = vault_kubernetes_auth_backend_role.opentelemetry_issuer.role_name
            secretRef = {
              name = kubernetes_secret.opentelemetry_issuer_token.metadata[0].name
              key  = "token"
            }
          }
        }
      }
    }
  }
}
