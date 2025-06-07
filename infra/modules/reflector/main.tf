provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

resource "kubernetes_namespace" "reflector" {
  metadata {
    name = var.namespace
  }
}

module "helm" {
  source    = "./modules/helm"
  namespace = kubernetes_namespace.reflector.metadata.0.name
  providers = {
    helm = helm
  }
}
