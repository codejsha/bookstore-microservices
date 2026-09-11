include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "kubernetes" {
  path = find_in_parent_folders("_envcommon/provider_kubernetes.hcl")
}

include "helm" {
  path = find_in_parent_folders("_envcommon/provider_helm.hcl")
}

remote_state {
  backend = "local"

  generate = {
    path      = "backend.tf"
    if_exists = "overwrite_terragrunt"
  }

  config = {
    path = "${get_terragrunt_dir()}/terraform.tfstate"
  }
}

inputs = {
  namespace               = "vault"
  vault_url               = "https://vault.example.com"
  vault_address           = "vault.example.com"
  vault_service_name      = "vault-ui"
  kube_api_server_address = "workstation.internal:6443"
  local_output_dir        = "${get_repo_root()}/infra/service/vault"
}

prevent_destroy = true

terraform {
  source = "${get_repo_root()}/infra//service/vault"
}
