terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
  }
}

resource "vault_pki_secret_backend_role" "harbor" {
  name    = "harbor"
  backend = "pki_int"
  allowed_domains = [
    "localhost",
    "harbor.example.com",
    "harbor.harbor.svc",
    "harbor.harbor.svc.cluster.local"
  ]
  allow_bare_domains = true
  allow_glob_domains = true
  max_ttl            = "1440h"
}

resource "kubernetes_service_account" "harbor_issuer" {
  metadata {
    name      = "harbor-issuer"
    namespace = var.namespace
  }
}

resource "kubernetes_secret" "harbor_issuer_token" {
  type = "kubernetes.io/service-account-token"
  metadata {
    name      = "harbor-issuer-token"
    namespace = var.namespace
    annotations = {
      "kubernetes.io/service-account.name" = kubernetes_service_account.harbor_issuer.metadata.0.name
    }
  }
}

resource "vault_kubernetes_auth_backend_role" "harbor_issuer" {
  role_name = "harbor-issuer"
  backend   = "kubernetes"
  bound_service_account_names = [kubernetes_service_account.harbor_issuer.metadata.0.name]
  bound_service_account_namespaces = [var.namespace]
  token_policies = ["pki_int"]
  token_ttl = "3600" # 1 hour
}

resource "kubernetes_cluster_role_binding" "harbor_issuer_token_rolebinding" {
  metadata {
    name = "harbor-issuer-token-rolebinding"
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = "system:auth-delegator"
  }
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.harbor_issuer.metadata.0.name
    namespace = var.namespace
  }
}

resource "kubernetes_manifest" "harbor_issuer" {
  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "Issuer"
    metadata = {
      name      = "harbor-issuer"
      namespace = var.namespace
    }
    spec = {
      vault = {
        server = "http://vault.vault.svc.cluster.local:8200"
        path   = "pki_int/sign/harbor"
        caBundle = base64encode(trimspace(file(var.kube_ca_crt_path)))
        auth = {
          kubernetes = {
            mountPath = "/v1/auth/kubernetes"
            role      = vault_kubernetes_auth_backend_role.harbor_issuer.role_name
            secretRef = {
              name = kubernetes_secret.harbor_issuer_token.metadata.0.name
              key  = "token"
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_manifest" "harbor_cert" {
  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "Certificate"
    metadata = {
      name      = "harbor-cert"
      namespace = var.namespace
    }
    spec = {
      secretName = "harbor-cert"
      commonName = var.harbor_address
      dnsNames = [
        "localhost",
        "harbor.example.com",
        "harbor.harbor.svc",
        "harbor.harbor.svc.cluster.local"
      ]
      issuerRef = {
        group = "cert-manager.io"
        kind  = "Issuer"
        name  = kubernetes_manifest.harbor_issuer.manifest.metadata.name
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
