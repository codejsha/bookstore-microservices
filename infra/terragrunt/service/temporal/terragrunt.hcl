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
    "../../service/postgresql",
    "../../service/gateway-api",
    "../../service/prometheus-crds",
    "../../service/istio",
  ]
}

inputs = {
  namespace        = "temporal"
  temporal_address = "temporal.example.com"

  temporal_authorization_enabled = true

  postgres_cluster_name       = "temporal-postgres"
  postgres_instances          = 1
  postgres_image              = "ghcr.io/cloudnative-pg/postgresql:18"
  postgres_storage_size       = "10Gi"
  postgres_storage_class_name = "local-path"

  postgres_resources = {
    requests = { cpu = "50m", memory = "512Mi" }
    limits   = { cpu = "1", memory = "1Gi" }
  }
}

prevent_destroy = true

terraform {
  source = "${get_repo_root()}/infra//service/temporal"
}
