terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "~> 5.10"
    }
    harbor = {
      source  = "goharbor/harbor"
      version = "~> 3.12"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.9"
    }
  }
}

ephemeral "vault_kv_secret_v2" "harbor_admin" {
  mount = "kv"
  name  = "harbor/admin/credentials"
}

provider "harbor" {
  url      = var.harbor_url
  username = ephemeral.vault_kv_secret_v2.harbor_admin.data["username"]
  password = ephemeral.vault_kv_secret_v2.harbor_admin.data["password"]
}

resource "random_password" "harbor_user" {
  for_each = toset(var.harbor_usernames)
  length   = 24
  special  = false
}

locals {
  harbor_users_set = toset([
    for u in var.harbor_usernames : {
      username = u
      email    = "${u}@example.com"
      password = random_password.harbor_user[u].result
    }
  ])
}

resource "vault_kv_secret_v2" "harbor_user" {
  for_each = toset(var.harbor_usernames)
  mount    = "kv"
  name     = "harbor/users/${each.value}/credentials"
  data_json = jsonencode({
    username = each.value
    password = random_password.harbor_user[each.value].result
  })
}

module "user" {
  source       = "./modules/user"
  harbor_users = local.harbor_users_set
  providers = {
    harbor = harbor
  }
}

module "project" {
  source          = "./modules/project"
  harbor_projects = var.harbor_projects
  providers = {
    harbor = harbor
  }
}

module "robot" {
  source     = "./modules/robot"
  robot_name = var.robot_name
  projects   = [for k, v in var.harbor_projects : v.project_name]
  depends_on = [module.project]
  providers = {
    harbor = harbor
    vault  = vault
  }
}
