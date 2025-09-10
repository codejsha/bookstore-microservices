terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
    vault = {
      source  = "hashicorp/vault"
      version = "~> 5.10"
    }
  }
}

resource "kubernetes_namespace_v1" "bookstore" {
  metadata {
    name   = var.namespace
    labels = var.namespace_labels
  }
}

resource "kubernetes_limit_range_v1" "resource_limits" {
  metadata {
    name      = "resource-limits"
    namespace = kubernetes_namespace_v1.bookstore.metadata[0].name
  }
  spec {
    limit {
      type = "Container"
      default_request = {
        cpu    = "10m"
        memory = "64Mi"
      }
      default = {
        cpu    = "500m"
        memory = "512Mi"
      }
      max = {
        cpu    = "2"
        memory = "2Gi"
      }
    }
  }
}

data "vault_kv_secret_v2" "harbor_pull" {
  mount = "kv"
  name  = var.harbor_pull_vault_path
}

resource "kubernetes_secret_v1" "harbor_pull" {
  metadata {
    name      = var.harbor_pull_secret_name
    namespace = kubernetes_namespace_v1.bookstore.metadata[0].name
  }

  type = "kubernetes.io/dockerconfigjson"

  data = {
    ".dockerconfigjson" = jsonencode({
      auths = {
        (var.registry_host) = {
          username = data.vault_kv_secret_v2.harbor_pull.data["username"]
          password = data.vault_kv_secret_v2.harbor_pull.data["password"]
          auth = base64encode(format("%s:%s",
            data.vault_kv_secret_v2.harbor_pull.data["username"],
            data.vault_kv_secret_v2.harbor_pull.data["password"],
          ))
        }
      }
    })
  }
}
