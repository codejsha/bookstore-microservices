terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "keycloak" {
  namespace  = var.namespace
  name       = "keycloak"
  repository = "oci://registry-1.docker.io/bitnamicharts"
  chart      = "keycloak"
  version    = "24.4.9"
  values = [
    file("${path.module}/values.yaml")
  ]
  timeout = 180
  set {
    name  = "auth.adminUser"
    value = var.admin_username
  }
  set {
    name  = "auth.adminPassword"
    value = var.admin_password
  }
}
