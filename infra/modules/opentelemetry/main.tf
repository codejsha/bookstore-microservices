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

resource "helm_release" "opentelemetry" {
  namespace  = kubernetes_namespace.opentelemetry.metadata.0.name
  name       = "opentelemetry"
  repository = "https://open-telemetry.github.io/opentelemetry-helm-charts"
  chart      = "opentelemetry-kube-stack"
  version    = "0.6.1"
  values = [
    file("${path.module}/values.yaml")
  ]
}
