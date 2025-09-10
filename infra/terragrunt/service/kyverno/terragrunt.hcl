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
    "../../service/cert-manager",
  ]
}

inputs = {
  namespace     = "kyverno"
  chart_version = "3.8.2"

  image_pull_secret   = "harbor-pull"
  ca_bundle_configmap = "internal-ca"
  ca_bundle_key       = "ca-bundle.crt"
}

terraform {
  source = "${get_repo_root()}/infra//service/kyverno"
}
