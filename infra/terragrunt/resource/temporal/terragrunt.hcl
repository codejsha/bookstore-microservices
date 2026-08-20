include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "kubernetes" {
  path = find_in_parent_folders("_envcommon/provider_kubernetes.hcl")
}

dependencies {
  paths = [
    "../../service/temporal",
  ]
}

inputs = {
  namespace           = "temporal"
  temporal_address    = "temporal-frontend.temporal.svc.cluster.local:7233"
  temporal_namespaces = ["bookstore"]
  retention           = "72h"
}

terraform {
  source = "${get_repo_root()}/infra//resource/temporal"
}
