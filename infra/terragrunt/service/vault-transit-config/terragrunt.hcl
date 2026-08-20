include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "kubernetes" {
  path = find_in_parent_folders("_envcommon/provider_kubernetes.hcl")
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
  main_vault_namespace = "vault"

  vault_unsealer_url   = get_env("VAULT_UNSEALER_ADDR", "http://127.0.0.1:8200")
  vault_unsealer_token = get_env("VAULT_UNSEALER_TOKEN", "")
}

prevent_destroy = true

terraform {
  source = "${get_repo_root()}/infra//service/vault-transit-config"
}
