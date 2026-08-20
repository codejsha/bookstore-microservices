terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "loki" {
  namespace  = var.namespace
  name       = "loki"
  repository = "https://grafana.github.io/helm-charts"
  chart      = "loki"
  version    = "6.16.0"
  values = [
    templatefile("${path.module}/values.yaml", {
      s3_endpoint = var.s3_endpoint
      s3_region   = var.s3_region
      s3_bucket   = var.s3_bucket
    })
  ]
  set_sensitive = [
    { name = "loki.storage.s3.accessKeyId", value = var.s3_access_key_id },
    { name = "loki.storage.s3.secretAccessKey", value = var.s3_secret_key }
  ]
}
