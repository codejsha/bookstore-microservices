terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "kubernetes_service_account_v1" "vault" {
  metadata {
    name      = "vault"
    namespace = var.namespace
  }
}

resource "kubernetes_secret_v1" "vault" {
  type = "kubernetes.io/service-account-token"
  metadata {
    name      = "vault-token"
    namespace = var.namespace
    annotations = {
      "kubernetes.io/service-account.name" = kubernetes_service_account_v1.vault.metadata[0].name
    }
  }
}

resource "helm_release" "vault" {
  depends_on = [kubernetes_service_account_v1.vault, kubernetes_secret_v1.vault]
  namespace  = var.namespace
  name       = "vault"
  repository = "https://helm.releases.hashicorp.com"
  chart      = "vault"
  version    = "0.32.0"
  values = [
    file("${path.module}/values.yaml")
  ]
  timeout = 120
}
