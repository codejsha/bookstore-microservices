terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "gitea" {
  namespace  = var.namespace
  name       = "gitea"
  repository = "oci://registry-1.docker.io/bitnamicharts"
  chart      = "gitea"
  version    = "3.1.10"
  values = [
    file("${path.module}/values.yaml")
  ]
  timeout = 120

  set {
    name  = "adminEmail"
    value = var.admin_email
  }
  set_sensitive {
    name  = "adminUsername"
    value = var.admin_username
  }
  set_sensitive {
    name  = "adminPassword"
    value = var.admin_password
  }
}
