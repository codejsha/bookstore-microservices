terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "cloudnative_pg" {
  namespace  = var.namespace
  name       = "cloudnative-pg"
  repository = "https://cloudnative-pg.github.io/charts"
  chart      = "cloudnative-pg"
  version    = "0.22.1"
  values = [
    file("${path.module}/operator-values.yaml"),
  ]
  timeout = 300
}
