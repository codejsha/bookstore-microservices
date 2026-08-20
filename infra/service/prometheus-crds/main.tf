terraform {
  required_providers {
    helm = {
      source  = "hashicorp/helm"
      version = "~> 3.2"
    }
  }
}

resource "helm_release" "prometheus_operator_crds" {
  namespace  = "kube-system"
  name       = "prometheus-operator-crds"
  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "prometheus-operator-crds"
  version    = "25.0.0"
}
