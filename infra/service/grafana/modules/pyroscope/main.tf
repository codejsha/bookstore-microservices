terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "pyroscope" {
  namespace  = var.namespace
  name       = "pyroscope"
  repository = "https://grafana.github.io/helm-charts"
  chart      = "pyroscope"
  version    = "1.14.0"
  values = [
    templatefile("${path.module}/values.yaml", {
      s3_endpoint = replace(replace(var.s3_endpoint, "http://", ""), "https://", "")
      s3_region   = var.s3_region
      s3_bucket   = var.s3_bucket
    })
  ]
  set_sensitive = [
    { name = "pyroscope.structuredConfig.storage.s3.access_key_id", value = var.s3_access_key_id },
    { name = "pyroscope.structuredConfig.storage.s3.secret_access_key", value = var.s3_secret_key }
  ]
}
