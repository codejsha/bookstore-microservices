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

resource "kubernetes_namespace_v1" "vso" {
  metadata {
    name = var.namespace
  }
}

module "helm" {
  source        = "./modules/helm"
  namespace     = kubernetes_namespace_v1.vso.metadata[0].name
  vault_address = var.vault_address
  providers = {
    helm = helm
  }
}
