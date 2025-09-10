terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "config_server" {
  namespace          = var.namespace
  name               = "config-server"
  repository         = "oci://harbor.example.com/bookstore-helm-charts"
  repository_ca_file = var.repository_ca_file
  chart              = "config-server"
  version            = var.chart_version
  values = [
    file("${path.module}/values.yaml")
  ]
}
