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

dependency "vault" {
  config_path = "../vault"

  mock_outputs = {
    approle_credentials = {
      seaweedfs = {
        role_id   = "mock-role-id"
        secret_id = "mock-secret-id"
      }
    }
    seaweedfs_s3_credentials = {
      access_key_id     = "mock-access-key-id"
      secret_access_key = "mock-secret-access-key"
    }
  }
  mock_outputs_allowed_terraform_commands = ["validate", "plan"]
}

inputs = {
  vault_role_id   = dependency.vault.outputs.approle_credentials["seaweedfs"].role_id
  vault_secret_id = dependency.vault.outputs.approle_credentials["seaweedfs"].secret_id
  s3_access_key   = dependency.vault.outputs.seaweedfs_s3_credentials.access_key_id
  s3_secret_key   = dependency.vault.outputs.seaweedfs_s3_credentials.secret_access_key

  namespace                  = "seaweedfs"
  seaweedfs_address          = "seaweedfs.example.com"
  seaweedfs_service_name     = "seaweedfs-admin"
  seaweedfs_api_address      = "seaweedfs-api.example.com"
  seaweedfs_api_service_name = "seaweedfs-s3"
  vault_url                  = "https://vault.example.com"
  s3_endpoint                = "https://seaweedfs-api.example.com"
  bucket_names               = ["terraform-state", "loki", "tempo", "pyroscope", "flink-checkpoints"]
}

prevent_destroy = true

terraform {
  source = "${get_repo_root()}/infra//service/seaweedfs"
}
