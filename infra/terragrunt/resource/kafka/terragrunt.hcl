include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "kubernetes" {
  path = find_in_parent_folders("_envcommon/provider_kubernetes.hcl")
}

include "helm" {
  path = find_in_parent_folders("_envcommon/provider_helm.hcl")
}

dependencies {
  paths = [
    "../../service/kafka",
  ]
}

inputs = {
  namespace = "kafka"
}

terraform {
  source = "${get_repo_root()}/infra//resource/kafka"
}
