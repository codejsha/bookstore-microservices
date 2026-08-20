terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "~> 5.10"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.3"
    }
    external = {
      source  = "hashicorp/external"
      version = "~> 2.4"
    }
  }
}

data "vault_kv_secret_v2" "gitea_admin" {
  mount = "kv"
  name  = "gitea/admin/credentials"
}

module "cloud_config" {
  source         = "./modules/cloud-config"
  org_name       = var.org_name
  gitea_api_url  = var.gitea_api_url
  gitea_http_url = var.gitea_http_url
  repo_name      = var.cloud_config_repo
  content_dir    = var.cloud_config_dir
  gitea_username = data.vault_kv_secret_v2.gitea_admin.data["username"]
  gitea_password = data.vault_kv_secret_v2.gitea_admin.data["password"]
}

module "helm_chart" {
  source         = "./modules/helm-chart"
  app_repos      = var.app_repos
  helm_dir       = var.helm_dir
  org_name       = var.org_name
  gitea_api_url  = var.gitea_api_url
  gitea_http_url = var.gitea_http_url
  gitea_username = data.vault_kv_secret_v2.gitea_admin.data["username"]
  gitea_password = data.vault_kv_secret_v2.gitea_admin.data["password"]
}

module "source_repo" {
  source         = "./modules/source-repo"
  services       = var.services
  repo_root      = var.repo_root
  repo_suffix    = var.source_repo_suffix
  org_name       = var.org_name
  gitea_api_url  = var.gitea_api_url
  gitea_http_url = var.gitea_http_url
  gitea_username = data.vault_kv_secret_v2.gitea_admin.data["username"]
  gitea_password = data.vault_kv_secret_v2.gitea_admin.data["password"]
}

module "github_mirror" {
  source         = "./modules/github-mirror"
  repos          = var.lib_repos
  github_owner   = var.github_owner
  org_name       = var.org_name
  gitea_api_url  = var.gitea_api_url
  gitea_http_url = var.gitea_http_url
  gitea_username = data.vault_kv_secret_v2.gitea_admin.data["username"]
  gitea_password = data.vault_kv_secret_v2.gitea_admin.data["password"]
}

module "tekton_catalog" {
  source         = "./modules/tekton-catalog"
  org_name       = var.org_name
  gitea_api_url  = var.gitea_api_url
  gitea_http_url = var.gitea_http_url
  repo_name      = var.tekton_catalog_repo
  upstream_url   = var.tekton_catalog_upstream
  gitea_username = data.vault_kv_secret_v2.gitea_admin.data["username"]
  gitea_password = data.vault_kv_secret_v2.gitea_admin.data["password"]
}

module "tekton_custom_catalog" {
  source         = "./modules/tekton-custom-catalog"
  org_name       = var.org_name
  gitea_api_url  = var.gitea_api_url
  gitea_http_url = var.gitea_http_url
  repo_name      = var.tekton_custom_catalog_repo
  content_dir    = var.tekton_custom_catalog_dir
  gitea_username = data.vault_kv_secret_v2.gitea_admin.data["username"]
  gitea_password = data.vault_kv_secret_v2.gitea_admin.data["password"]
}
