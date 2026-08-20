include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "kubernetes" {
  path = find_in_parent_folders("_envcommon/provider_kubernetes.hcl")
}

include "helm" {
  path = find_in_parent_folders("_envcommon/provider_helm.hcl")
}

dependencies {
  paths = [
    "../bookstore-base",
  ]
}

inputs = {
  namespace = "bookstore"

  chart_version    = "3.0.31"
  image_registry   = "docker.io"
  image_repository = "bitnamilegacy/valkey"
  image_tag        = "8.1.3-debian-12-r3"

  valkey_services = {
    catalog = {
      storage_size = "2Gi"
      memory_limit = "256Mi"
    }
    identity = {
      storage_size = "2Gi"
      memory_limit = "256Mi"
    }
    oauth2-proxy = {
      storage_size = "2Gi"
      memory_limit = "256Mi"
    }
    order = {
      storage_size = "2Gi"
      memory_limit = "256Mi"
    }
    ratelimit = {
      storage_size = "2Gi"
      memory_limit = "256Mi"
    }
    payment = {
      storage_size = "2Gi"
      memory_limit = "256Mi"
    }
  }
}

terraform {
  source = "${get_repo_root()}/infra//resource/valkey"
}
