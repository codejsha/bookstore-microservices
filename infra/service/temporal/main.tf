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
    vault = {
      source  = "hashicorp/vault"
      version = "~> 5.10"
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

data "vault_kv_secret_v2" "mysql" {
  mount = "kv"
  name  = "temporal/mysql"
}

module "mysql" {
  source             = "./modules/mysql"
  namespace          = kubernetes_namespace_v1.temporal.metadata[0].name
  root_password      = data.vault_kv_secret_v2.mysql.data["root_password"]
  cluster_name       = var.mysql_cluster_name
  instances          = var.mysql_instances
  router_instances   = var.mysql_router_instances
  mysql_version      = var.mysql_version
  storage_size       = var.mysql_storage_size
  storage_class_name = var.mysql_storage_class_name
  db_secret_name     = var.mysql_db_secret_name
  mysql_resources    = var.mysql_resources
  router_resources   = var.mysql_router_resources
  providers = {
    kubernetes = kubernetes
  }
}

module "temporal" {
  source          = "./modules/helm"
  namespace       = kubernetes_namespace_v1.temporal.metadata[0].name
  db_connect_addr = module.mysql.connect_addr
  db_secret_name  = module.mysql.db_secret_name
  db_user         = var.mysql_db_user
  providers = {
    helm = helm
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
