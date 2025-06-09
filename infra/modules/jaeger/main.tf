provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

resource "kubernetes_namespace" "jaeger" {
  metadata {
    name = var.namespace
    labels = {
      "istio-injection" = "enabled"
    }
  }
}

module "helm" {
  source    = "./modules/helm"
  namespace = kubernetes_namespace.jaeger.metadata.0.name
}

module "istio" {
  source       = "./modules/istio"
  namespace    = kubernetes_namespace.jaeger.metadata.0.name
  host_address = var.jaeger_address
  host_fqdn    = var.jaeger_fqdn
  dest_port    = 16686
  name_prefix  = "jaeger"
}
