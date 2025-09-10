terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "kyverno" {
  namespace  = var.namespace
  name       = "kyverno"
  repository = "https://kyverno.github.io/kyverno/"
  chart      = "kyverno"
  version    = var.chart_version
  timeout    = 600
  values = [
    templatefile("${path.module}/values.yaml", {
      image_pull_secret   = var.image_pull_secret
      ca_bundle_configmap = var.ca_bundle_configmap
      ca_bundle_key       = var.ca_bundle_key
    })
  ]
}
