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
  repository = "https://dl.gitea.com/charts"
  chart      = "gitea"
  version    = "12.1.0"
  values = [
    file("${path.module}/values.yaml")
  ]
  timeout = 180

  set_sensitive {
    name  = "gitea.admin.email"
    value = var.admin_email
  }
  set_sensitive {
    name  = "gitea.admin.username"
    value = var.admin_username
  }
  set_sensitive {
    name  = "gitea.admin.password"
    value = var.admin_password
  }
  set_sensitive {
    name  = "valkey.global.valkey.password"
    value = var.valkey_password
  }
  set_sensitive {
    name  = "postgresql.global.postgresql.auth.username"
    value = var.admin_username
  }
  set_sensitive {
    name  = "postgresql.global.postgresql.auth.password"
    value = var.admin_password
  }
}
