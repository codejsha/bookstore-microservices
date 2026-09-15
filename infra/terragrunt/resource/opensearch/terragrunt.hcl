include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "kubernetes" {
  path = find_in_parent_folders("_envcommon/provider_kubernetes.hcl")
}

include "vault" {
  path = find_in_parent_folders("_envcommon/provider_vault.hcl")
}

dependencies {
  paths = [
    "../../service/opensearch",
    "../../service/vault",
    "../../resource/cert-manager-issuers",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role = "opensearch-config"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  opensearch_api_url = "https://opensearch-api.example.com"

  keycloak_issuer_url = "https://keycloak.example.com/realms/platform-infra"
  oidc_client_id      = "opensearch-dashboards"
}

terraform {
  source = "${get_repo_root()}/infra//resource/opensearch"
}
