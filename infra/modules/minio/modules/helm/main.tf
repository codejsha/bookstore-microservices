terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "kubernetes_service_account" "minio" {
  metadata {
    name      = "minio"
    namespace = var.namespace
  }
}

resource "kubernetes_secret" "minio" {
  type = "kubernetes.io/service-account-token"
  metadata {
    name      = "minio-token"
    namespace = var.namespace
    annotations = {
      "kubernetes.io/service-account.name" = kubernetes_service_account.minio.metadata[0].name
    }
  }
}

resource "helm_release" "minio" {
  depends_on = [kubernetes_service_account.minio, kubernetes_secret.minio]
  namespace  = var.namespace
  name       = "minio"
  repository = "oci://registry-1.docker.io/bitnamicharts"
  chart      = "minio"
  version    = "17.0.1"
  values = [
    file("${path.module}/values.yaml")
  ]
  timeout = 120
}
