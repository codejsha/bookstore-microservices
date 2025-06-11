terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "opentelemetry" {
  namespace  = var.namespace
  name       = "opentelemetry"
  repository = "https://open-telemetry.github.io/opentelemetry-helm-charts"
  chart      = "opentelemetry-kube-stack"
  version    = "0.6.1"
  values = [
    file("${path.module}/values.yaml")
  ]
  timeout = 180

  set_sensitive {
    name  = "collectors.cluster.config.extensions.basicauth/opensearch.client_auth.username"
    value = var.opensearch_username
  }
  set_sensitive {
    name  = "collectors.cluster.config.extensions.basicauth/opensearch.client_auth.password"
    value = var.opensearch_password
  }
}
