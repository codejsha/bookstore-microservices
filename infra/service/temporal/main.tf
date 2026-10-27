terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 3.2"
    }
  }
}

resource "kubernetes_namespace_v1" "temporal" {
  metadata {
    name = var.namespace
    labels = {
      "istio.io/dataplane-mode" = "ambient"
    }
  }
}

resource "kubernetes_limit_range_v1" "resource_limits" {
  metadata {
    name      = "resource-limits"
    namespace = kubernetes_namespace_v1.temporal.metadata[0].name
  }
  spec {
    limit {
      type = "Container"
      default_request = {
        cpu    = "10m"
        memory = "32Mi"
      }
    }
  }
}

module "postgresql" {
  source             = "./modules/postgresql"
  namespace          = kubernetes_namespace_v1.temporal.metadata[0].name
  cluster_name       = var.postgres_cluster_name
  instances          = var.postgres_instances
  image              = var.postgres_image
  storage_size       = var.postgres_storage_size
  storage_class_name = var.postgres_storage_class_name
  resources          = var.postgres_resources
  providers = {
    kubernetes = kubernetes
  }
}

module "temporal" {
  source          = "./modules/helm"
  namespace       = kubernetes_namespace_v1.temporal.metadata[0].name
  db_connect_addr = module.postgresql.connect_addr
  db_secret_name  = module.postgresql.db_secret_name
  db_user         = module.postgresql.db_user
  depends_on      = [module.postgresql]

  authorization_enabled = var.temporal_authorization_enabled
  providers = {
    helm = helm
  }
}

resource "kubernetes_manifest" "internal_frontend_authz" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "temporal-internal-frontend-authz"
      namespace = kubernetes_namespace_v1.temporal.metadata[0].name
    }
    spec = {
      selector = {
        matchLabels = {
          "app.kubernetes.io/instance"  = "temporal"
          "app.kubernetes.io/component" = "internal-frontend"
        }
      }
      action = "ALLOW"
      rules = [
        {
          from = [
            { source = { namespaces = [kubernetes_namespace_v1.temporal.metadata[0].name] } }
          ]
        },
        {
          to = [
            { operation = { ports = ["9090"] } }
          ]
        }
      ]
    }
  }
}

module "gateway_api" {
  source          = "../../shared/gateway-api"
  hostname        = var.temporal_address
  service_name    = "temporal-web"
  service_port    = 8080
  route_namespace = kubernetes_namespace_v1.temporal.metadata[0].name
  name_prefix     = "temporal"
  request_timeout = "10m"
}
