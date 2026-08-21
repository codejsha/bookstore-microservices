terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.37.1"
    }
    vault = {
      source  = "hashicorp/vault"
      version = ">= 5.0.0"
    }
  }
}

provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "vault" {
  address      = var.vault_url
  token        = var.vault_token
  ca_cert_file = var.kube_ca_crt_path
}

resource "kubernetes_namespace" "bookstore" {
  metadata {
    name = var.namespace
    labels = {
      "istio-injection" = "enabled"
    }
  }
}

module "service_account" {
  source         = "./modules/service-account"
  namespace      = kubernetes_namespace.bookstore.metadata[0].name
  mysql_services = var.mysql_services
}

module "secret" {
  source                 = "./modules/secret"
  namespace              = kubernetes_namespace.bookstore.metadata[0].name
  mysql_services         = var.mysql_services
  mysql_service_accounts = module.service_account.service_accounts
  mysql_root_password    = var.mysql_root_password
  mysql_username         = var.mysql_username
  mysql_password         = var.mysql_password
  providers = {
    vault = vault
  }
}
