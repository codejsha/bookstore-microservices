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

resource "kubernetes_namespace_v1" "cnpg_system" {
  metadata {
    name = var.namespace
    labels = {
      "istio.io/dataplane-mode" = "ambient"
    }
  }
}

module "postgresql_operator" {
  source    = "./modules/postgresql-operator"
  namespace = kubernetes_namespace_v1.cnpg_system.metadata[0].name
  providers = {
    helm = helm
  }
}
