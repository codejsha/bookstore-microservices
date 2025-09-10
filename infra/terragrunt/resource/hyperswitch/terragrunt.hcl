include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "vault" {
  path = find_in_parent_folders("_envcommon/provider_vault.hcl")
}

dependencies {
  paths = [
    "../../service/vault",
    "../../service/hyperswitch",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role = "hyperswitch-config"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  hyperswitch_admin_url = "https://hyperswitch-api.example.com"

  merchant_id   = "bookstore"
  merchant_name = "Bookstore"
}

terraform {
  source = "${get_repo_root()}/infra//resource/hyperswitch"
}
