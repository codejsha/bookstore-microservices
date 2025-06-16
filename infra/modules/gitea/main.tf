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
  username = var.admin_username
  password = var.admin_password
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
  namespace      = kubernetes_namespace.gitea.metadata[0].name
  admin_email    = var.admin_email
  admin_username = var.admin_username
  admin_password = var.admin_password
  providers = {
    helm = helm
  }
}

module "istio" {
  source       = "./modules/istio"
  namespace    = kubernetes_namespace.gitea.metadata[0].name
  host_address = var.gitea_address
  host_fqdn    = var.gitea_fqdn
  dest_port    = 3000
  name_prefix  = "gitea"
}

module "organization" {
  source   = "./modules/organization"
  org_name = var.org_name
  providers = {
    gitea = gitea
  }
}

module "repo_ssh" {
  source    = "./modules/repo-ssh"
  namespace = kubernetes_namespace.gitea.metadata[0].name
  providers = {
    vault = vault
  }
}

module "repos" {
  for_each = toset(concat(var.dev_repos, var.devops_repos))
  source           = "./modules/repo"
  gitea_username   = var.org_name
  repo_name        = each.key
  repo_description = "${each.key} repository"
  providers = {
    gitea = gitea
    vault = vault
  }
}

module "dev_team" {
  source           = "./modules/team"
  org_name         = var.org_name
  team_name        = "dev-team"
  user_repos       = var.dev_repos
  user_credentials = var.dev_user_credentials
  providers = {
    gitea = gitea
  }
}

module "devops_team" {
  source           = "./modules/team"
  org_name         = var.org_name
  team_name        = "devops-team"
  user_repos       = var.devops_repos
  user_credentials = var.devops_user_credentials
  providers = {
    gitea = gitea
  }
}

module "token" {
  source = "./modules/token"
  providers = {
    gitea = gitea
    vault = vault
  }
}

module "cert" {
  source           = "./modules/cert"
  namespace        = kubernetes_namespace.gitea.metadata[0].name
  gitea_address    = var.gitea_address
  kube_ca_crt_path = var.kube_ca_crt_path
  providers = {
    vault = vault
  }
}
