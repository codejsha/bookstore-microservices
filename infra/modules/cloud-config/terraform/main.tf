provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

resource "kubernetes_namespace" "config_server" {
  metadata {
    name = var.namespace
    labels = {
      "istio-injection" = "enabled"
    }
  }
}

module "helm" {
  source             = "./modules/helm"
  namespace          = kubernetes_namespace.config_server.metadata.0.name
  repository_ca_file = var.repository_ca_file
  providers = {
    helm = helm
  }
}
