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
  temporal_address    = "temporal-internal-frontend.temporal.svc.cluster.local:7236"
  temporal_namespaces = ["bookstore"]
  retention           = "72h"
}

terraform {
  source = "${get_repo_root()}/infra//resource/temporal"
}
