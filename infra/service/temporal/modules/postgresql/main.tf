terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

resource "kubernetes_manifest" "cluster" {
  manifest = {
    apiVersion = "postgresql.cnpg.io/v1"
    kind       = "Cluster"
    metadata = {
      name      = var.cluster_name
      namespace = var.namespace
    }
    spec = {
      instances             = var.instances
      imageName             = var.image
      primaryUpdateStrategy = "unsupervised"
      enableSuperuserAccess = true
      resources             = var.resources
      storage = merge(
        { size = var.storage_size },
        var.storage_class_name != "" ? { storageClass = var.storage_class_name } : {},
      )
      monitoring = {
        enablePodMonitor = true
      }
      inheritedMetadata = {
        labels = {
          release = "prometheus"
        }
      }
    }
  }
  computed_fields = ["metadata.labels", "metadata.annotations", "spec"]
}
