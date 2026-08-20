terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

resource "kubernetes_secret_v1" "cluster" {
  metadata {
    name      = "${var.cluster_name}-cluster-secret"
    namespace = var.namespace
  }
  type = "Opaque"
  data = {
    rootUser     = "root"
    rootHost     = "%"
    rootPassword = var.root_password
  }
}

resource "kubernetes_secret_v1" "db" {
  metadata {
    name      = var.db_secret_name
    namespace = var.namespace
  }
  type = "Opaque"
  data = {
    password = var.root_password
  }
}

resource "kubernetes_manifest" "innodb_cluster" {
  manifest = {
    apiVersion = "mysql.oracle.com/v2"
    kind       = "InnoDBCluster"
    metadata = {
      name      = var.cluster_name
      namespace = var.namespace
    }
    spec = {
      instances        = var.instances
      secretName       = kubernetes_secret_v1.cluster.metadata[0].name
      tlsUseSelfSigned = true
      version          = var.mysql_version
      router = {
        instances = var.router_instances
        podSpec = {
          containers = [
            {
              name      = "router"
              resources = var.router_resources
            }
          ]
        }
      }
      datadirVolumeClaimTemplate = {
        accessModes      = ["ReadWriteOnce"]
        storageClassName = var.storage_class_name
        resources = {
          requests = {
            storage = var.storage_size
          }
        }
      }
      podSpec = {
        containers = [
          {
            name      = "mysql"
            resources = var.mysql_resources
          }
        ]
      }
    }
  }
  computed_fields = ["metadata.labels", "metadata.annotations", "spec"]
}
