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

resource "vault_pki_secret_backend_role" "youtrack" {
  name    = "youtrack"
  backend = "pki_int"
  allowed_domains = [
    "localhost",
    var.youtrack_address,
    "youtrack.${var.namespace}.svc",
    "youtrack.${var.namespace}.svc.cluster.local",
  ]
  allow_bare_domains = true
  allow_glob_domains = true
  max_ttl            = "1440h"
}

resource "kubernetes_service_account_v1" "youtrack_issuer" {
  metadata {
    name      = "youtrack-issuer"
    namespace = var.namespace
  }
}

resource "kubernetes_secret_v1" "youtrack_issuer_token" {
  type = "kubernetes.io/service-account-token"
  metadata {
    name      = "youtrack-issuer-token"
    namespace = var.namespace
    annotations = {
      "kubernetes.io/service-account.name" = kubernetes_service_account_v1.youtrack_issuer.metadata[0].name
    }
  }
}

resource "vault_kubernetes_auth_backend_role" "youtrack_issuer" {
  role_name                        = "youtrack-issuer"
  backend                          = "kubernetes"
  bound_service_account_names      = [kubernetes_service_account_v1.youtrack_issuer.metadata[0].name]
  bound_service_account_namespaces = [var.namespace]
  token_policies                   = ["pki_int"]
  token_ttl                        = "3600"
}

resource "kubernetes_cluster_role_binding_v1" "youtrack_issuer_token_rolebinding" {
  metadata {
    name = "youtrack-issuer-token-rolebinding"
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = "system:auth-delegator"
  }
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account_v1.youtrack_issuer.metadata[0].name
    namespace = var.namespace
  }
}

resource "kubernetes_manifest" "youtrack_issuer" {
  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "Issuer"
    metadata = {
      name      = "youtrack-issuer"
      namespace = var.namespace
    }
    spec = {
      vault = {
        server   = "http://vault.vault.svc.cluster.local:8200"
        path     = "pki_int/sign/youtrack"
        caBundle = base64encode(trimspace(var.kube_ca_cert))
        auth = {
          kubernetes = {
            mountPath = "/v1/auth/kubernetes"
            role      = vault_kubernetes_auth_backend_role.youtrack_issuer.role_name
            secretRef = {
              name = kubernetes_secret_v1.youtrack_issuer_token.metadata[0].name
              key  = "token"
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_manifest" "youtrack_cert" {
  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "Certificate"
    metadata = {
      name      = "youtrack-cert"
      namespace = var.namespace
    }
    spec = {
      secretName = "youtrack-cert"
      commonName = var.youtrack_address
      dnsNames = [
        "localhost",
        var.youtrack_address,
        "youtrack.${var.namespace}.svc",
        "youtrack.${var.namespace}.svc.cluster.local",
      ]
      issuerRef = {
        group = "cert-manager.io"
        kind  = "Issuer"
        name  = kubernetes_manifest.youtrack_issuer.manifest.metadata.name
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
