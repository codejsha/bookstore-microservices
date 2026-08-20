terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "~> 5.10"
    }
    nexus = {
      source  = "datadrivers/nexus"
      version = "~> 3.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.9"
    }
  }
}

ephemeral "vault_kv_secret_v2" "nexus_admin" {
  mount = "kv"
  name  = "nexus/admin/credentials"
}

provider "nexus" {
  url      = var.nexus_url
  username = ephemeral.vault_kv_secret_v2.nexus_admin.data["username"]
  password = ephemeral.vault_kv_secret_v2.nexus_admin.data["password"]
  insecure = false
}

resource "random_password" "publisher" {
  length  = 24
  special = false
}

module "repository" {
  source                    = "./modules/repository"
  proxy_name                = var.proxy_name
  hosted_name               = var.hosted_name
  group_name                = var.group_name
  maven_proxy_name          = var.maven_proxy_name
  maven_hosted_name         = var.maven_hosted_name
  maven_group_name          = var.maven_group_name
  blob_store_name           = var.blob_store_name
  npm_registry_remote_url   = var.npm_registry_remote_url
  maven_registry_remote_url = var.maven_registry_remote_url
  maven_plugins_proxy_name  = var.maven_plugins_proxy_name
  maven_plugins_remote_url  = var.maven_plugins_remote_url
  raw_hosted_name           = var.raw_hosted_name
  pypi_proxy_name           = var.pypi_proxy_name
  pypi_group_name           = var.pypi_group_name
  pypi_registry_remote_url  = var.pypi_registry_remote_url
  hosted_write_policy       = var.hosted_write_policy

  go_proxy_name          = var.go_proxy_name
  go_group_name          = var.go_group_name
  go_registry_remote_url = var.go_registry_remote_url

  docker_proxy_name          = var.docker_proxy_name
  docker_registry_remote_url = var.docker_registry_remote_url
  docker_http_port           = var.docker_http_port
  providers = {
    nexus = nexus
  }
}

module "security" {
  source                       = "./modules/security"
  role_id                      = var.publisher_role_id
  publisher_username           = var.publisher_username
  publisher_password           = random_password.publisher.result
  hosted_repository_name       = module.repository.hosted_name
  group_repository_name        = module.repository.group_name
  maven_hosted_repository_name = module.repository.maven_hosted_name
  maven_group_repository_name  = module.repository.maven_group_name
  pypi_group_repository_name   = module.repository.pypi_group_name
  raw_hosted_repository_name   = module.repository.raw_hosted_name
  go_group_repository_name     = module.repository.go_group_name
  docker_proxy_repository_name = module.repository.docker_proxy_name
  providers = {
    nexus = nexus
  }
}

resource "vault_kv_secret_v2" "publisher" {
  mount = "kv"
  name  = "nexus/publisher"
  data_json = jsonencode({
    username = var.publisher_username
    password = random_password.publisher.result
  })
}

output "npm_group_url" {
  description = "Group repository URL — point .npmrc registry here (resolves hosted then proxy)."
  value       = "${var.nexus_url}/repository/${var.group_name}/"
}

output "npm_hosted_url" {
  description = "Hosted repository URL — publish target for private @bookstore/* packages."
  value       = "${var.nexus_url}/repository/${var.hosted_name}/"
}

output "maven_group_url" {
  description = "Maven group URL — what the Gradle builds resolve against (hosted first, then the cached Maven Central)."
  value       = "${var.nexus_url}/repository/${var.maven_group_name}/"
}

output "maven_hosted_url" {
  description = "Maven hosted URL — publish target for the private com.codejsha.platform artifacts."
  value       = "${var.nexus_url}/repository/${var.maven_hosted_name}/"
}

output "pypi_group_url" {
  description = "PyPI group URL — uv/pip index is <this>/simple."
  value       = "${var.nexus_url}/repository/${var.pypi_group_name}/"
}

output "npm_proxy_url" {
  description = "Proxy repository URL — caches registry.npmjs.org."
  value       = "${var.nexus_url}/repository/${var.proxy_name}/"
}

output "go_group_url" {
  description = "Go group URL — what the Go builds' GOPROXY resolves against (caches proxy.golang.org)."
  value       = "${var.nexus_url}/repository/${var.go_group_name}/"
}

output "docker_proxy_endpoint" {
  description = <<-EOT
    host:port CI prefixes onto an image reference to pull it through the ghcr.io
    cache. Not a /repository/ URL — docker clients address the connector port
    directly, so this is the in-cluster Service fronting that connector.
  EOT
  value       = "${var.docker_service_name}.${var.namespace}.svc.cluster.local:${var.docker_http_port}"
}
