terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "config_server" {
  namespace  = var.namespace
  name       = "config-server"
  repository = "oci://harbor.example.com/bookstore"
  chart      = "config-server"
  version    = "0.1.0"
  repository_ca_file = file(var.repository_ca_file)
  values = [
    file("${path.module}/values.yaml")
  ]
}
