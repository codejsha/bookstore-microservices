include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
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

terraform {
  source = "${get_repo_root()}/infra//service/gateway-api"
}
