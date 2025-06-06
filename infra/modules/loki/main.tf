provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

resource "kubernetes_namespace" "loki" {
  metadata {
    name = var.namespace
    labels = {
      "istio-injection" = "enabled"
    }
  }
}

module "helm" {
  source         = "./modules/helm"
  namespace      = kubernetes_namespace.loki.metadata.0.name
  minio_username = var.minio_username
  minio_password = var.minio_password
  providers = {
    helm = helm
  }
}
