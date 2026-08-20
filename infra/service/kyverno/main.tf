terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 3.2"
    }
  }
}

resource "kubernetes_namespace_v1" "kyverno" {
  metadata {
    name = var.namespace
  }
}

module "helm" {
  source              = "./modules/helm"
  namespace           = kubernetes_namespace_v1.kyverno.metadata[0].name
  chart_version       = var.chart_version
  image_pull_secret   = var.image_pull_secret
  ca_bundle_configmap = var.ca_bundle_configmap
  ca_bundle_key       = var.ca_bundle_key
  providers = {
    helm = helm
  }
}
