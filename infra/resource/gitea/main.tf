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
    gitea = {
      source  = "go-gitea/gitea"
      version = "~> 0.8"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.9"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.3"
    }
  }
}

resource "random_password" "dev" {
  for_each = toset(var.dev_usernames)
  length   = 20
  special  = false
}

resource "random_password" "devops" {
  for_each = toset(var.devops_usernames)
  length   = 20
  special  = false
}

resource "random_password" "webhook" {
  length  = 40
  special = false
}

resource "vault_kv_secret_v2" "webhook" {
  mount     = "kv"
  name      = "gitea/webhook/credentials"
  data_json = jsonencode({ token = random_password.webhook.result })
}

resource "tls_private_key" "tekton_resolver" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "gitea_public_key" "tekton_resolver" {
  username  = var.admin_username
  title     = "tekton-resolver"
  key       = trimspace(tls_private_key.tekton_resolver.public_key_openssh)
  read_only = true
}

resource "vault_kv_secret_v2" "tekton_resolver_ssh" {
  mount = "kv"
  name  = "gitea/ssh/tekton-resolver"
  data_json = jsonencode({
    private = tls_private_key.tekton_resolver.private_key_pem
    public  = tls_private_key.tekton_resolver.public_key_openssh
  })
}

locals {
  dev_user_credentials    = [for u in var.dev_usernames : { username = u, password = random_password.dev[u].result }]
  devops_user_credentials = [for u in var.devops_usernames : { username = u, password = random_password.devops[u].result }]
}

ephemeral "vault_kv_secret_v2" "gitea_admin" {
  mount = "kv"
  name  = "gitea/admin/credentials"
}

provider "gitea" {
  base_url = var.gitea_url
  username = var.admin_username
  password = ephemeral.vault_kv_secret_v2.gitea_admin.data["password"]
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

module "config_server_access" {
  source          = "./modules/config-server-access"
  namespace       = "cloud-config"
  service_account = "config-server"
  providers = {
    vault = vault
  }
}

module "repos" {
  depends_on       = [module.organization]
  for_each         = toset(concat(var.dev_repos, var.devops_repos))
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
  depends_on       = [module.organization, module.repos]
  source           = "./modules/team"
  org_name         = var.org_name
  team_name        = "dev-team"
  user_repos       = var.dev_repos
  user_credentials = local.dev_user_credentials
  providers = {
    gitea = gitea
  }
}

module "devops_team" {
  depends_on       = [module.organization, module.repos]
  source           = "./modules/team"
  org_name         = var.org_name
  team_name        = "devops-team"
  user_repos       = var.devops_repos
  user_credentials = local.devops_user_credentials
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
