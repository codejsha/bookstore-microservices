include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "kubernetes" {
  path = find_in_parent_folders("_envcommon/provider_kubernetes.hcl")
}

include "helm" {
  path = find_in_parent_folders("_envcommon/provider_helm.hcl")
}

include "vault" {
  path = find_in_parent_folders("_envcommon/provider_vault.hcl")
}

dependencies {
  paths = [
    "../../service/vault",
    "../../service/mysql",
    "../../service/gateway-api",
    "../../service/prometheus-crds",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role = "temporal"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  namespace        = "temporal"
  temporal_address = "temporal.example.com"

  mysql_cluster_name       = "temporal-mysql"
  mysql_instances          = 1
  mysql_router_instances   = 1
  mysql_version            = "8.4.3"
  mysql_storage_size       = "10Gi"
  mysql_storage_class_name = "local-path"
  mysql_db_secret_name     = "temporal-db-secret"
  mysql_db_user            = "root"

  mysql_resources = {
    requests = { cpu = "50m", memory = "1280Mi" }
    limits   = { cpu = "1", memory = "3Gi" }
  }
  mysql_router_resources = {
    requests = { cpu = "10m", memory = "128Mi" }
    limits   = { cpu = "500m", memory = "256Mi" }
  }
}

prevent_destroy = true

terraform {
  source = "${get_repo_root()}/infra//service/temporal"
}
