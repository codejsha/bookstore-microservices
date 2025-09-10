terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "harbor" {
  namespace  = var.namespace
  name       = "harbor"
  repository = "https://helm.goharbor.io"
  chart      = "harbor"
  version    = "1.18.2"
  values = [
    file("${path.module}/values.yaml"),
  ]
  set = [
    { name = "existingSecretAdminPassword", value = "harbor-admin-credentials" },
    { name = "existingSecretAdminPasswordKey", value = "password" },
    { name = "registry.credentials.existingSecret", value = "harbor-registry-credentials" },
    { name = "persistence.imageChartStorage.s3.bucket", value = var.s3_bucket },
    { name = "persistence.imageChartStorage.s3.regionendpoint", value = var.s3_endpoint },
  ]
  set_sensitive = [
    { name = "persistence.imageChartStorage.s3.accesskey", value = var.aws_access_key },
    { name = "persistence.imageChartStorage.s3.secretkey", value = var.aws_secret_key },
  ]
}
