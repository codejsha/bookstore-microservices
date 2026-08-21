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
    argocd = {
      source  = "argoproj-labs/argocd"
      version = ">= 7.8.2"
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

provider "argocd" {
  server_addr = "${var.argocd_address}:80"
  username    = var.argocd_username
  password    = var.argocd_password
  plain_text  = true
}

module "token" {
  source = "./modules/token"
  providers = {
    argocd = argocd
    vault  = vault
  }
}

module "repo" {
  source     = "./modules/repo"
  org_name   = var.org_name
  gitea_fqdn = var.gitea_fqdn
  app_repos  = var.app_repos
  providers = {
    argocd = argocd
    vault  = vault
  }
}

module "project" {
  source     = "./modules/project"
  namespace  = var.namespace
  org_name   = var.org_name
  gitea_fqdn = var.gitea_fqdn
  app_repos  = var.app_repos
  providers = {
    argocd = argocd
  }
}
