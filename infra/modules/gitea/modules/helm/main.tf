terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "kubernetes_service_account" "gitea" {
  metadata {
    name      = "gitea"
    namespace = var.namespace
  }
}

resource "kubernetes_secret" "gitea" {
  type = "kubernetes.io/service-account-token"
  metadata {
    name      = "gitea-token"
    namespace = var.namespace
    annotations = {
      "kubernetes.io/service-account.name" = kubernetes_service_account.gitea.metadata[0].name
    }
  }
}

resource "helm_release" "gitea" {
  namespace  = var.namespace
  name       = "gitea"
  repository = "oci://registry-1.docker.io/bitnamicharts"
  chart      = "gitea"
  version    = "3.1.10"
  values = [
    file("${path.module}/values.yaml")
  ]
  timeout = 120

  set_sensitive {
    name  = "adminEmail"
    value = var.admin_email
  }
  set_sensitive {
    name  = "adminUsername"
    value = var.admin_username
  }
  set_sensitive {
    name  = "adminPassword"
    value = var.admin_password
  }
}
