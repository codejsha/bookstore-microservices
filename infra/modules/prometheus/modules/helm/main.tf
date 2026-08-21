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
  version    = "73.0.0"
  values = [
    file("${path.module}/values.yaml")
  ]

  set {
    name  = "grafana.additionalDataSources[0].basicAuth"
    value = true
  }
  set_sensitive {
    name  = "grafana.additionalDataSources[0].basicAuthUser"
    value = var.opensearch_username
  }
  set_sensitive {
    name  = "grafana.additionalDataSources[0].secureJsonData.basicAuthPassword"
    value = var.opensearch_password
  }
}
