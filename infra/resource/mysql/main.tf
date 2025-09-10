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

resource "random_password" "mysql_root" {
  length  = 24
  special = false
}

resource "random_password" "mysql_service" {
  for_each = toset(var.mysql_services)
  length   = 24
  special  = false
}

locals {
  mysql_credentials = [
    for svc in var.mysql_services : {
      service  = svc
      username = svc
      password = random_password.mysql_service[svc].result
    }
  ]
}

module "service_account" {
  source         = "./modules/service-account"
  namespace      = var.namespace
  mysql_services = var.mysql_services
}

module "secret" {
  source                 = "./modules/secret"
  namespace              = var.namespace
  mysql_services         = var.mysql_services
  mysql_service_accounts = module.service_account.service_accounts
  mysql_root_password    = random_password.mysql_root.result
  mysql_credentials      = local.mysql_credentials
  providers = {
    vault = vault
  }
}

module "cluster" {
  source                 = "./modules/cluster"
  namespace              = var.namespace
  mysql_services         = var.mysql_services
  mysql_db_config        = var.mysql_db_config
  mysql_client_image     = var.mysql_client_image
  mysqld_exporter_image  = var.mysqld_exporter_image
  cross_service_readonly = var.cross_service_readonly
  providers = {
    kubernetes = kubernetes
  }
  depends_on = [module.secret]
}

module "db_engine" {
  source                 = "./modules/db-engine"
  namespace              = var.namespace
  mysql_services         = var.mysql_services
  mysql_db_config        = var.mysql_db_config
  cross_service_readonly = var.cross_service_readonly
  providers = {
    vault      = vault
    kubernetes = kubernetes
  }
  depends_on = [module.cluster]
}
