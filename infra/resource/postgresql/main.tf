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
    random = {
      source  = "hashicorp/random"
      version = "~> 3.9"
    }
  }
}

resource "random_password" "postgres_service" {
  for_each = toset(var.postgres_services)
  length   = 24
  special  = false
}

resource "random_password" "debezium" {
  length  = 24
  special = false
}

locals {
  postgres_credentials = [
    for svc in var.postgres_services : {
      service  = svc
      username = svc
      password = random_password.postgres_service[svc].result
    }
  ]
}

module "service_account" {
  source            = "./modules/service-account"
  namespace         = var.namespace
  postgres_services = var.postgres_services
}

module "secret" {
  source                    = "./modules/secret"
  namespace                 = var.namespace
  postgres_services         = var.postgres_services
  postgres_service_accounts = module.service_account.service_accounts
  postgres_credentials      = local.postgres_credentials
  debezium_username         = var.debezium_username
  debezium_password         = random_password.debezium.result
  providers = {
    vault = vault
  }
}

module "cluster" {
  source    = "./modules/cluster"
  namespace = var.namespace
  clusters  = var.postgres_cluster_config

  depends_on = [module.secret]

  providers = {
    kubernetes = kubernetes
  }
}

module "db_engine" {
  source = "./modules/db-engine"

  namespace             = var.namespace
  postgres_services     = var.postgres_services
  postgres_db_map       = var.postgres_db_map
  postgres_app_user_map = var.postgres_app_user_map

  depends_on = [module.cluster]

  providers = {
    vault      = vault
    kubernetes = kubernetes
  }
}
