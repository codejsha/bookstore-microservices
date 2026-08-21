terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
  }
}

resource "vault_pki_secret_backend_role" "keycloak" {
  name    = "keycloak"
  backend = "pki_int"
  allowed_domains = [
    "localhost",
    "keycloak.example.com",
    "keycloak.keycloak.svc",
    "keycloak.keycloak.svc.cluster.local"
  ]
  allow_bare_domains = true
  allow_glob_domains = true
  max_ttl            = "1440h"
}

resource "kubernetes_service_account" "keycloak_issuer" {
  metadata {
    name      = "keycloak-issuer"
    namespace = var.namespace
  }
}

resource "kubernetes_secret" "keycloak_issuer_secret" {
  metadata {
    name      = "keycloak-issuer-secret"
    namespace = var.namespace
    annotations = {
      "kubernetes.io/service-account.name" : kubernetes_service_account.keycloak_issuer.metadata.0.name
    }
  }
  type = "kubernetes.io/service-account-token"
}

resource "vault_kubernetes_auth_backend_role" "keycloak_issuer" {
  role_name = "keycloak-issuer"
  backend   = "kubernetes"
  bound_service_account_names = [kubernetes_service_account.keycloak_issuer.metadata.0.name]
  bound_service_account_namespaces = [var.namespace]
  token_policies = ["pki_int"]
}

resource "kubernetes_cluster_role_binding" "keycloak_issuer_token_rolebinding" {
  metadata {
    name = "keycloak-issuer-token-rolebinding"
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = "system:auth-delegator"
  }
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.keycloak_issuer.metadata.0.name
    namespace = var.namespace
  }
}

resource "kubernetes_manifest" "keycloak_issuer" {
  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "Issuer"
    metadata = {
      name      = "keycloak-issuer"
      namespace = var.namespace
    }
    spec = {
      vault = {
        server = "http://vault.vault.svc.cluster.local:8200"
        path   = "pki_int/sign/keycloak"
        caBundle = base64encode(data.kubernetes_secret.vault_tls.data["kubernetes-ca.crt"])
        auth = {
          kubernetes = {
            mountPath = "/v1/auth/kubernetes"
            role      = vault_kubernetes_auth_backend_role.keycloak_issuer.role_name
            secretRef = {
              name = kubernetes_secret.keycloak_issuer_secret.metadata.0.name
              key : "token"
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_manifest" "keycloak_cert" {
  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "Certificate"
    metadata = {
      name      = "keycloak-cert"
      namespace = var.namespace
    }
    spec = {
      secretName = "keycloak-tls"
      commonName = var.keycloak_address
      dnsNames = [
        "localhost",
        "keycloak.example.com",
        "keycloak.keycloak.svc",
        "keycloak.keycloak.svc.cluster.local"
      ]
      issuerRef = {
        group = "cert-manager.io"
        kind  = "Issuer"
        name  = kubernetes_manifest.keycloak_issuer.manifest.metadata.name
      }
    }
  }
}
