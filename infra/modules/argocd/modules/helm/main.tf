terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "argocd" {
  namespace  = var.namespace
  name       = "argocd"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = "7.8.2"
  values = [
    file("${path.module}/values.yaml")
  ]
}
