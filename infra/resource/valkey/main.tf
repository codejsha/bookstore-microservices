terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 3.2"
    }
  }
}

resource "helm_release" "valkey" {
  for_each = var.valkey_services

  namespace  = var.namespace
  name       = "${each.key}-valkey"
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "valkey"
  version    = var.chart_version

  values = [yamlencode({
    global = {
      security = {
        allowInsecureImages = true
      }
    }
    image = {
      registry   = var.image_registry
      repository = var.image_repository
      tag        = var.image_tag
    }
    architecture = "standalone"
    auth = {
      enabled = false
    }
    networkPolicy = {
      enabled = false
    }
    metrics = {
      enabled = true
      image = {
        registry   = var.image_registry
        repository = "bitnamilegacy/redis-exporter"
        tag        = "1.76.0-debian-12-r0"
      }
      serviceMonitor = {
        enabled = true
        additionalLabels = {
          release = "prometheus"
        }
      }
      resourcesPreset = ""
      resources = {
        requests = { cpu = "10m", memory = "32Mi" }
        limits   = { cpu = "100m", memory = "64Mi" }
      }
    }
    primary = {
      resourcesPreset = ""
      resources = {
        requests = { cpu = "10m", memory = "64Mi" }
        limits   = { cpu = "200m", memory = each.value.memory_limit }
      }
      persistence = {
        enabled = true
        size    = each.value.storage_size
      }
    }
  })]
}
