terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "loki" {
  namespace  = var.namespace
  name       = "loki"
  repository = "https://grafana.github.io/helm-charts"
  chart      = "loki"
  version    = "6.30.1"
  values = [
    file("${path.module}/values.yaml"),
  ]

  set_sensitive {
    name  = "loki.storageConfig.s3.accessKeyId"
    value = var.minio_username
  }
  set_sensitive {
    name  = "loki.storageConfig.s3.secretAccessKey"
    value = var.minio_password
  }
}
