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
    "../../service/tekton",
    "../../service/vault",
    "../../resource/argocd",
    "../../resource/cluster-dns",
    "../../resource/gitea",
    "../../resource/gitea-repo",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  gitea      = read_terragrunt_config(find_in_parent_folders("_envcommon/gitea.hcl")).locals
  vault_role = "tekton-config"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  max_concurrent_pipelineruns = 10

  namespace      = "bookstore-ci"
  gitea_address  = "git.example.com"
  gitea_org      = local.gitea.gitea_org
  harbor_address = "harbor.example.com"

  nexus_maven_url = "http://nexus.nexus.svc.cluster.local:8081/repository/maven-group/"
  nexus_npm_url   = "http://nexus.nexus.svc.cluster.local:8081/repository/npm-group/"
  nexus_pypi_url  = "http://nexus.nexus.svc.cluster.local:8081/repository/pypi-group/simple"
  nexus_raw_url   = "http://nexus.nexus.svc.cluster.local:8081/repository/raw-hosted"

  nexus_docker_registry = "nexus-docker.nexus.svc.cluster.local:5000"

  codegen_verify_targets = []

  source_repo_suffix = local.gitea.source_repo_suffix

  gradle_services = ["order", "payment", "settlement", "admin"]
  golang_services = ["catalog", "customer", "identity", "inventory"]
  uv_services     = ["delivery", "notification", "support"]
  vite_services   = ["web", "admin-web"]

  service_accounts = ["git-source", "git-deploy"]
}

terraform {
  source = "${get_repo_root()}/infra//resource/tekton"
}
