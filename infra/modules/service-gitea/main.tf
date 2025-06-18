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
    gitea = {
      source  = "go-gitea/gitea"
      version = ">= 0.6.0"
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

provider "gitea" {
  base_url = var.gitea_url
  username = var.admin_username
  password = var.admin_password
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
  namespace = var.namespace
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
