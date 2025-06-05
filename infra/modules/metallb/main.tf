provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

resource "kubernetes_namespace" "metallb" {
  metadata {
    name = var.namespace
  }
}

# module "legacy" {
#   source         = "./modules/legacy"
#   namespace      = kubernetes_namespace.metallb.metadata.0.name
#   pool_addresses = var.pool_addresses
#   providers = {
#     kubernetes = kubernetes
#     helm       = helm
#   }
# }

module "helm" {
  source    = "./modules/helm"
  namespace = kubernetes_namespace.metallb.metadata.0.name
  providers = {
    helm = helm
  }
}

module "pool" {
  source         = "./modules/pool"
  namespace      = kubernetes_namespace.metallb.metadata.0.name
  pool_addresses = var.pool_addresses
  providers = {
    kubernetes = kubernetes
  }
}
