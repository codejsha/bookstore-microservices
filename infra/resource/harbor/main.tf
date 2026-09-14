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
  mount = "kv-infra"
  name  = "harbor/admin/credentials"
}

provider "harbor" {
  url      = var.harbor_url
  username = ephemeral.vault_kv_secret_v2.harbor_admin.data["username"]
  password = ephemeral.vault_kv_secret_v2.harbor_admin.data["password"]
}

data "vault_kv_secret_v2" "harbor_oidc" {
  mount = "kv-infra"
  name  = "keycloak/harbor-oidc/client-secret"
}

module "auth" {
  source             = "./modules/auth"
  oidc_endpoint      = var.oidc_issuer
  oidc_client_id     = data.vault_kv_secret_v2.harbor_oidc.data["client_id"]
  oidc_client_secret = data.vault_kv_secret_v2.harbor_oidc.data["client_secret"]
  providers = {
    harbor = harbor
  }
}

module "project" {
  source          = "./modules/project"
  harbor_projects = var.harbor_projects
  depends_on      = [module.auth]
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
