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

provider "vault" {
  address = var.vault_unsealer_url
  token   = var.vault_unsealer_token
}

module "transit_unseal" {
  source               = "./modules/transit-unseal"
  main_vault_namespace = var.main_vault_namespace
  providers = {
    vault      = vault
    kubernetes = kubernetes
  }
}
