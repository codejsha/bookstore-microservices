terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "tempo" {
  namespace  = var.namespace
  name       = "tempo"
  repository = "https://grafana.github.io/helm-charts"
  chart      = "tempo"
  version    = "1.23.0"
  values = [
    templatefile("${path.module}/values.yaml", {
      s3_endpoint = replace(replace(var.s3_endpoint, "http://", ""), "https://", "")
      s3_region   = var.s3_region
      s3_bucket   = var.s3_bucket
    })
  ]
  set_sensitive = [
    { name = "tempo.storage.trace.s3.access_key", value = var.s3_access_key_id },
    { name = "tempo.storage.trace.s3.secret_key", value = var.s3_secret_key }
  ]
}
