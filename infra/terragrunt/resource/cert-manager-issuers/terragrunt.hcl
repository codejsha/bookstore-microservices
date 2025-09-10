include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "kubernetes" {
  path = find_in_parent_folders("_envcommon/provider_kubernetes.hcl")
}

dependencies {
  paths = [
    "../../service/cert-manager",
    "../../service/vault",
  ]
}

locals {
  common = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
}

inputs = {
  vault_url = local.common.vault_url
}

terraform {
  source = "${get_repo_root()}/infra//resource/cert-manager-issuers"
}
