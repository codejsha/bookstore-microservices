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

resource "kubernetes_namespace_v1" "vault_transit" {
  metadata {
    name = var.namespace
    labels = {
      "istio.io/dataplane-mode" = "ambient"
    }
  }
}

module "helm" {
  source    = "./modules/helm"
  namespace = kubernetes_namespace_v1.vault_transit.metadata[0].name
  providers = {
    helm = helm
  }
}
