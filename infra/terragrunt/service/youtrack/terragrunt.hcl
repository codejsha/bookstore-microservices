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
    "../../service/gateway-api",
    "../../service/cert-manager",
    "../../service/grafana",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role = "youtrack"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  namespace        = "youtrack"
  youtrack_address = "youtrack.example.com"
  admin_email      = "admin@example.com"

  grafana_url = "https://grafana.example.com"

  chart_version    = "3.4.0"
  youtrack_version = "2026.1.13570"

  storage_class       = "local-path"
  data_storage_size   = "50Gi"
  logs_storage_size   = "5Gi"
  backup_storage_size = "50Gi"
}

terraform {
  source = "${get_repo_root()}/infra//service/youtrack"
}
