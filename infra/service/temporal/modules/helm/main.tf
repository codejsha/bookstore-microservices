terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "temporal" {
  namespace  = var.namespace
  name       = "temporal"
  repository = "https://go.temporal.io/helm-charts"
  chart      = "temporal"
  version    = "1.1.1"
  values = [
    templatefile("${path.module}/values.yaml", {
      db_connect_addr = var.db_connect_addr
      db_secret_name  = var.db_secret_name
      db_user         = var.db_user
    })
  ]
  timeout = 900
}
