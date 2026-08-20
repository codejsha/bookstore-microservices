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
    "../../service/prometheus-crds",
    "../../service/vault-secrets-operator",
    "../../service/grafana",
  ]
}

locals {
  common         = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role     = "gitea"
  admin_username = "gitea_admin"
  admin_email    = "admin@example.com"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  namespace          = "gitea"
  gitea_address      = "git.example.com"
  gitea_service_name = "gitea-http"
  admin_username     = local.admin_username
  admin_email        = local.admin_email

  grafana_url = "https://grafana.example.com"
}

terraform {
  source = "${get_repo_root()}/infra//service/gitea"
}
