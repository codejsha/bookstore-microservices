terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
  }
}

resource "vault_pki_secret_backend_role" "gitea" {
  name    = "gitea"
  backend = "pki_int"
  allowed_domains = [
    "localhost",
    "git.example.com",
    "gitea.gitea.svc",
    "gitea.gitea.svc.cluster.local"
  ]
  allow_bare_domains = true
  allow_glob_domains = true
  max_ttl            = "1440h"
}

resource "kubernetes_service_account" "gitea_issuer" {
  metadata {
    name      = "gitea-issuer"
    namespace = var.namespace
  }
}

resource "kubernetes_secret" "gitea_issuer_token" {
  type = "kubernetes.io/service-account-token"
  metadata {
    name      = "gitea-issuer-token"
    namespace = var.namespace
    annotations = {
      "kubernetes.io/service-account.name" = kubernetes_service_account.gitea_issuer.metadata[0].name
    }
  }
}

resource "vault_kubernetes_auth_backend_role" "gitea_issuer" {
  role_name = "gitea-issuer"
  backend   = "kubernetes"
  bound_service_account_names = [kubernetes_service_account.gitea_issuer.metadata[0].name]
  bound_service_account_namespaces = [var.namespace]
  token_policies = ["pki_int"]
  token_ttl = "3600" # 1 hour
}

resource "kubernetes_cluster_role_binding" "gitea_issuer_token_rolebinding" {
  metadata {
    name = "gitea-issuer-token-rolebinding"
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = "system:auth-delegator"
  }
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.gitea_issuer.metadata[0].name
    namespace = var.namespace
  }
}

resource "kubernetes_manifest" "gitea_issuer" {
  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "Issuer"
    metadata = {
      name      = "gitea-issuer"
      namespace = var.namespace
    }
    spec = {
      vault = {
        server = "http://vault.vault.svc.cluster.local:8200"
        path   = "pki_int/sign/gitea"
        caBundle = base64encode(trimspace(file(var.kube_ca_crt_path)))
        auth = {
          kubernetes = {
            mountPath = "/v1/auth/kubernetes"
            role      = vault_kubernetes_auth_backend_role.gitea_issuer.role_name
            secretRef = {
              name = kubernetes_secret.gitea_issuer_token.metadata[0].name
              key  = "token"
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_manifest" "gitea_cert" {
  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "Certificate"
    metadata = {
      name      = "gitea-cert"
      namespace = var.namespace
    }
    spec = {
      secretName = "gitea-cert"
      commonName = var.gitea_address
      dnsNames = [
        "localhost",
        "git.example.com",
        "gitea.gitea.svc",
        "gitea.gitea.svc.cluster.local"
      ]
      issuerRef = {
        group = "cert-manager.io"
        kind  = "Issuer"
        name  = kubernetes_manifest.gitea_issuer.manifest.metadata.name
      }
      secretTemplate = {
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
