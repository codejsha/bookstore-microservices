terraform {
  required_providers {
    gitea = {
      source  = "go-gitea/gitea"
      version = "0.6.0"
    }
  }
}

provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

provider "vault" {
  address      = var.vault_url
  token        = var.vault_token
  ca_cert_file = var.kube_ca_crt_path
}

provider "gitea" {
  base_url = var.gitea_url
  username = var.gitea_username
  password = var.gitea_password
}

resource "kubernetes_namespace" "gitea" {
  metadata {
    name = var.namespace
    labels = {
      "istio-injection" = "enabled"
    }
  }
}

module "helm" {
  source         = "./modules/helm"
  namespace      = kubernetes_namespace.gitea.metadata.0.name
  admin_email    = var.admin_email
  admin_username = var.admin_username
  admin_password = var.admin_password
  providers = {
    helm = helm
  }
}

module "istio" {
  source       = "./modules/istio"
  namespace    = kubernetes_namespace.gitea.metadata.0.name
  host_address = var.gitea_address
  host_fqdn    = var.gitea_fqdn
  name_prefix  = "gitea"
}

module "ssh" {
  source    = "./modules/ssh"
  namespace = kubernetes_namespace.gitea.metadata.0.name
  providers = {
    vault = vault
  }
}

module "user" {
  source         = "./modules/user"
  admin_password = var.admin_password
  providers = {
    gitea = gitea
  }
}

module "organization" {
  source   = "./modules/organization"
  org_name = var.org_name
  providers = {
    gitea = gitea
    vault = vault
  }
}

module "repository" {
  for_each = toset(concat(var.dev_repos, var.devops_repos))
  source           = "./modules/repository"
  gitea_username   = var.org_name
  repo_name        = each.key
  repo_description = "${each.key} repository"
  providers = {
    gitea = gitea
  }
}

module "team" {
  source       = "./modules/team"
  org_name     = var.org_name
  dev_users    = var.dev_users
  dev_repos    = var.dev_repos
  devops_users = var.devops_users
  devops_repos = var.devops_repos
  providers = {
    gitea = gitea
    vault = vault
  }
}
