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

resource "helm_release" "flink_operator" {
  namespace  = var.namespace
  name       = "flink-operator"
  repository = "https://downloads.apache.org/flink/flink-kubernetes-operator-1.12.1/"
  chart      = "flink-kubernetes-operator"
  version    = "1.12.1"
  timeout    = 180
  values = [
    yamlencode({
      operatorPod = {
        resources = {
          requests = { cpu = "50m", memory = "512Mi" }
          limits   = { cpu = "500m", memory = "1536Mi" }
        }
      }
    })
  ]
}
