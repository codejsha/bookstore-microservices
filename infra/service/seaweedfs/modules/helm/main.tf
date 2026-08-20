terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "seaweedfs" {
  namespace  = var.namespace
  name       = "seaweedfs"
  repository = "https://seaweedfs.github.io/seaweedfs/helm"
  chart      = "seaweedfs"
  version    = "4.17.0"
  values = [
    file("${path.module}/values.yaml")
  ]
  set_sensitive = [
    { name = "s3.credentials.admin.accessKey", value = var.seaweedfs_access_key },
    { name = "s3.credentials.admin.secretKey", value = var.seaweedfs_secret_key },
    { name = "admin.secret.adminUser", value = var.seaweedfs_access_key },
    { name = "admin.secret.adminPassword", value = var.seaweedfs_secret_key }
  ]
  timeout = 600
}
