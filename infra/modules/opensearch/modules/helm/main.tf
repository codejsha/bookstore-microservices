terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "opensearch" {
  namespace  = var.namespace
  name       = "opensearch"
  repository = "https://opensearch-project.github.io/helm-charts"
  chart      = "opensearch"
  version    = "3.0.0"
  values = [
    file("${path.module}/values.yaml"),
  ]
  timeout = 180
}

resource "helm_release" "opensearch_dashboards" {
  depends_on = [helm_release.opensearch]
  namespace  = var.namespace
  name       = "opensearch-dashboards"
  repository = "https://opensearch-project.github.io/helm-charts"
  chart      = "opensearch-dashboards"
  version    = "3.0.0"
  values = [
    file("${path.module}/values-dashboards.yaml"),
  ]
  timeout = 180
}
