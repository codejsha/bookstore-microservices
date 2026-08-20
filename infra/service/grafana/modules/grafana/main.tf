terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "grafana" {
  namespace  = var.namespace
  name       = "grafana"
  repository = "https://grafana.github.io/helm-charts"
  chart      = "grafana"
  version    = "8.8.5"
  values = [
    file("${path.module}/values.yaml")
  ]
}
