terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "kubernetes_service_account_v1" "youtrack" {
  metadata {
    name      = "youtrack"
    namespace = var.namespace
  }
}

resource "kubernetes_secret_v1" "youtrack" {
  type = "kubernetes.io/service-account-token"
  metadata {
    name      = "youtrack-token"
    namespace = var.namespace
    annotations = {
      "kubernetes.io/service-account.name" = kubernetes_service_account_v1.youtrack.metadata[0].name
    }
  }
}

resource "helm_release" "youtrack" {
  namespace  = var.namespace
  name       = "youtrack"
  repository = "https://twenty-20.github.io/helm-charts/"
  chart      = "youtrack"
  version    = var.chart_version
  values = [
    file("${path.module}/values.yaml")
  ]
  timeout = 600
  set = [
    { name = "image.tag", value = var.youtrack_version },
    { name = "config.baseUrl", value = "https://${var.youtrack_address}" },
    { name = "ingress.enabled", value = "false" },
    { name = "persistence.data.storageClassName", value = var.storage_class },
    { name = "persistence.data.storageSize", value = var.data_storage_size },
    { name = "persistence.logs.storageClassName", value = var.storage_class },
    { name = "persistence.logs.storageSize", value = var.logs_storage_size },
    { name = "persistence.conf.storageClassName", value = var.storage_class },
    { name = "persistence.backups.volumeStorage.storageClassName", value = var.storage_class },
    { name = "persistence.backups.volumeStorage.storageSize", value = var.backup_storage_size },
  ]
}
