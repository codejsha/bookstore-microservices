terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

resource "kubernetes_service_account_v1" "nexus" {
  metadata {
    name      = "nexus"
    namespace = var.namespace
  }
}

resource "helm_release" "nexus" {
  namespace  = var.namespace
  name       = "nexus"
  repository = "https://sonatype.github.io/helm3-charts/"
  chart      = "nexus-repository-manager"
  version    = "64.2.0"
  values = [
    file("${path.module}/values.yaml")
  ]
  timeout = 600

  depends_on = [kubernetes_service_account_v1.nexus]
}

resource "kubernetes_service_v1" "docker" {
  metadata {
    name      = "nexus-docker"
    namespace = var.namespace
  }
  spec {
    selector = {
      "app.kubernetes.io/name"     = "nexus-repository-manager"
      "app.kubernetes.io/instance" = "nexus"
    }
    port {
      name        = "docker"
      port        = var.docker_port
      target_port = var.docker_port
    }
  }

  depends_on = [helm_release.nexus]
}
