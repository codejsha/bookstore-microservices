provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

resource "kubernetes_namespace" "istio-system" {
  metadata {
    name = var.namespace
  }
}

module "component" {
  source    = "./modules/component"
  namespace = kubernetes_namespace.istio-system.metadata.0.name
}

module "kiali" {
  source    = "./modules/kiali"
  namespace = kubernetes_namespace.istio-system.metadata.0.name
}

module "istio_kiali" {
  source       = "./modules/istio"
  namespace    = kubernetes_namespace.istio-system.metadata.0.name
  host_address = var.kiali_address
  host_fqdn    = var.kiali_fqdn
  dest_port    = 20001
  name_prefix  = "kiali"
}
