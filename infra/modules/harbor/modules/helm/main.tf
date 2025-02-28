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
  repository = "oci://registry-1.docker.io/bitnamicharts"
  chart      = "harbor"
  version    = "24.3.0"
  values = [
    file("${path.module}/values.yaml"),
  ]
  set {
    name  = "adminPassword"
    value = var.admin_password
  }
  set {
    name  = "persistence.imageChartStorage.s3.accesskey"
    value = var.aws_access_key
  }
  set {
    name  = "persistence.imageChartStorage.s3.secretkey"
    value = var.aws_secret_key
  }
}
