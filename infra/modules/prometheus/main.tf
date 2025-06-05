provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

resource "kubernetes_namespace" "prometheus" {
  metadata {
    name = var.namespace
    labels = {
      "istio-injection" = "enabled"
    }
  }
}

module "helm" {
  source    = "./modules/helm"
  namespace = kubernetes_namespace.prometheus.metadata.0.name
  providers = {
    helm = helm
  }
}

module istio_prometheus {
  source       = "./modules/istio"
  namespace    = kubernetes_namespace.prometheus.metadata.0.name
  host_address = var.prometheus_address
  host_fqdn    = var.prometheus_fqdn
  dest_port    = 9090
  name_prefix  = "prometheus"
}

module "istio_grafana" {
  source       = "./modules/istio"
  namespace    = kubernetes_namespace.prometheus.metadata.0.name
  host_address = var.grafana_address
  host_fqdn    = var.grafana_fqdn
  dest_port    = 80
  name_prefix  = "grafana"
}
