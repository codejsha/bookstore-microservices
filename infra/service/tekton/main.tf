terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 3.2"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.3"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.9"
    }
  }
}

module "component" {
  source = "./modules/component"
}

module "tektonconfig" {
  source                         = "./modules/tektonconfig"
  namespace                      = var.namespace
  prune_ttl_seconds              = var.prune_ttl_seconds
  prune_successful_history_limit = var.prune_successful_history_limit
  prune_failed_history_limit     = var.prune_failed_history_limit
  vault_internal_url             = var.vault_internal_url
  cosign_key_name                = var.cosign_key_name
  chains_vault_role              = var.chains_vault_role
  builder_id                     = var.builder_id

  gitea_internal_url    = var.gitea_internal_url
  gitea_org             = var.gitea_org
  resolver_token_secret = var.resolver_token_secret
  ca_bundle_configmap   = var.ca_bundle_configmap
  ca_bundle_key         = var.ca_bundle_key

  depends_on = [module.component]
}

resource "kubernetes_limit_range_v1" "resource_limits" {
  metadata {
    name      = "resource-limits"
    namespace = "tekton-pipelines"
  }
  spec {
    limit {
      type = "Container"
      default_request = {
        cpu    = "10m"
        memory = "32Mi"
      }
    }
  }
  depends_on = [module.component]
}

module "route" {
  source          = "../../shared/gateway-api"
  hostname        = var.tekton_address
  service_name    = var.tekton_service_name
  service_port    = 9097
  route_namespace = var.namespace
  name_prefix     = "tekton-dashboard"
  request_timeout = "10m"

  depends_on = [module.tektonconfig]
}
