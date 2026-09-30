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
  repository = "https://archive.apache.org/dist/flink/flink-kubernetes-operator-1.16.1/"
  chart      = "flink-kubernetes-operator"
  version    = "1.16.1"
  skip_crds  = true
  timeout    = 180
  values = [
    yamlencode({
      operatorPod = {
        resources = {
          requests = { cpu = "50m", memory = "512Mi" }
          limits   = { cpu = "500m", memory = "1536Mi" }
        }
        webhook = {
          resources = {
            requests = { cpu = "10m", memory = "256Mi" }
            limits   = { cpu = "500m", memory = "768Mi" }
          }
        }
      }
    })
  ]
}
