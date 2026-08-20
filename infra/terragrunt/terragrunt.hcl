locals {
  s3_bucket   = "terraform-state"
  s3_region   = "us-east-1"
  s3_endpoint = "https://seaweedfs-api.example.com"
}

remote_state {
  backend = "s3"

  generate = {
    path      = "backend.tf"
    if_exists = "overwrite_terragrunt"
  }

  config = {
    bucket = local.s3_bucket
    key    = "infra/${path_relative_to_include()}/terraform.tfstate"
    region = local.s3_region

    endpoints = {
      s3 = local.s3_endpoint
    }

    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_requesting_account_id  = true
    use_path_style              = true
  }
}

inputs = {}
