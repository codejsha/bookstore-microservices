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
    argocd = {
      source  = "argoproj-labs/argocd"
      version = "~> 7.16"
    }
    time = {
      source  = "hashicorp/time"
      version = "~> 0.14"
    }
  }
}

ephemeral "vault_kv_secret_v2" "argocd" {
  mount = "kv"
  name  = "argocd/admin/credentials"
}

provider "argocd" {
  server_addr = "${var.argocd_address}:443"
  username    = ephemeral.vault_kv_secret_v2.argocd.data["username"]
  password    = ephemeral.vault_kv_secret_v2.argocd.data["password"]
}

module "account_token" {
  source = "./modules/account-token"
  providers = {
    argocd = argocd
    vault  = vault
  }
}

module "repo" {
  source         = "./modules/repo"
  org_name       = var.org_name
  gitea_ssh_fqdn = var.gitea_ssh_fqdn
  app_repos      = var.app_repos
  providers = {
    argocd = argocd
    vault  = vault
  }
}

module "project" {
  source         = "./modules/project"
  namespace      = var.namespace
  org_name       = var.org_name
  gitea_ssh_fqdn = var.gitea_ssh_fqdn
  app_repos      = var.app_repos
  environment    = var.environment
  providers = {
    argocd = argocd
  }
}
