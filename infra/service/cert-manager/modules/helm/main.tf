terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "cert_manager" {
  namespace  = var.namespace
  name       = "cert-manager"
  repository = "https://charts.jetstack.io"
  chart      = "cert-manager"
  version    = "v1.20.0"
  set = [
    { name = "crds.enabled", value = true },
    { name = "prometheus.servicemonitor.enabled", value = true },
    { name = "prometheus.servicemonitor.labels.release", value = "prometheus" }
  ]
  values = [
    yamlencode({
      resources = {
        requests = { cpu = "10m", memory = "64Mi" }
        limits   = { cpu = "200m", memory = "256Mi" }
      }
      webhook = {
        resources = {
          requests = { cpu = "10m", memory = "64Mi" }
          limits   = { cpu = "200m", memory = "128Mi" }
        }
      }
      cainjector = {
        resources = {
          requests = { cpu = "10m", memory = "64Mi" }
          limits   = { cpu = "200m", memory = "256Mi" }
        }
      }
    })
  ]
}

resource "helm_release" "trust_manager" {
  namespace  = var.namespace
  name       = "trust-manager"
  repository = "https://charts.jetstack.io"
  chart      = "trust-manager"
  version    = "v0.20.1"
  values = [
    yamlencode({
      resources = {
        requests = { cpu = "10m", memory = "64Mi" }
        limits   = { cpu = "200m", memory = "128Mi" }
      }
      defaultPackage = {
        resources = {
          requests = { cpu = "10m", memory = "64Mi" }
          limits   = { cpu = "200m", memory = "128Mi" }
        }
      }
    })
  ]

  depends_on = [helm_release.cert_manager]
}
