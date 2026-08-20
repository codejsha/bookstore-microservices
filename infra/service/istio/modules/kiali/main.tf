terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "kiali" {
  namespace  = var.namespace
  name       = "kiali"
  repository = "https://kiali.org/helm-charts"
  chart      = "kiali-server"
  version    = var.kiali_chart_version
  values = [
    file("${path.module}/values.yaml")
  ]

  set_sensitive = [
    {
      name  = "external_services.grafana.auth.username"
      value = var.grafana_username
    },
    {
      name  = "external_services.grafana.auth.password"
      value = var.grafana_password
    }
  ]
}
