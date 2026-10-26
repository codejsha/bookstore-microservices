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
    "../bookstore-base",
    "../../service/mysql",
    "../../service/vault",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role = "mysql-config"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
    "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  namespace = "bookstore"
  mysql_services = []

  mysql_client_image = "mysql:8.4"

  mysql_db_config = {}

  cross_service_readonly = []
}

prevent_destroy = true

terraform {
  source = "${get_repo_root()}/infra//resource/mysql"
}
