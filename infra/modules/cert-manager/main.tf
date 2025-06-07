provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

resource "kubernetes_namespace" "cert_manager" {
  metadata {
    name = var.namespace
  }
}

module "helm" {
  source    = "./modules/helm"
  namespace = kubernetes_namespace.cert_manager.metadata.0.name
}

module "issuer" {
  source = "./modules/issuer"
}
