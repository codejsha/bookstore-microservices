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
  repository = "oci://harbor.example.com/main"
  chart      = "config-server"
  version    = "0.1.0"
}
