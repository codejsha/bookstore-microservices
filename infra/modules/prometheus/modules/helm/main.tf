terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "prometheus" {
  namespace  = var.namespace
  name       = "prometheus"
  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "kube-prometheus-stack"
  version    = "72.8.0"
  values = [
    file("${path.module}/values.yaml")
  ]
}
