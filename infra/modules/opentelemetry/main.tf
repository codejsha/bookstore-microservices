provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

resource "kubernetes_namespace" "opentelemetry" {
  metadata {
    name = var.namespace
    labels = {
      "istio-injection" = "enabled"
    }
  }
}

module "helm" {
  source    = "./modules/helm"
  namespace = kubernetes_namespace.opentelemetry.metadata.0.name
  providers = {
    helm = helm
  }
}

module "istio_collector" {
  source       = "./modules/istio"
  namespace    = kubernetes_namespace.opentelemetry.metadata.0.name
  name_prefix  = "otel-collector"
  host_address = var.otel_collector_address
  host_fqdn    = var.otel_collector_fqdn
}
