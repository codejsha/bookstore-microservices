terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.37.1"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "= 2.17.0"
    }
  }
}

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
  source              = "./modules/helm"
  namespace           = kubernetes_namespace.opentelemetry.metadata[0].name
  opensearch_username = var.opensearch_username
  opensearch_password = var.opensearch_password
  providers = {
    helm = helm
  }
}

module "istio" {
  source       = "./modules/istio"
  namespace    = kubernetes_namespace.opentelemetry.metadata[0].name
  host_address = var.otel_collector_http_address
  host_fqdn    = var.otel_collector_fqdn
  dest_port    = 4318
  name_prefix  = "otelcol-http"
}

module "istio_grpc" {
  source       = "./modules/istio-grpc"
  namespace    = kubernetes_namespace.opentelemetry.metadata[0].name
  host_address = var.otel_collector_grpc_address
  host_fqdn    = var.otel_collector_fqdn
  dest_port    = 4317
  name_prefix  = "otelcol-grpc"
}
