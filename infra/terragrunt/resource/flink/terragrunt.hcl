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
    "../../service/flink",
    "../../service/opensearch",
    "../../service/seaweedfs",
    "../../service/vault",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role = "flink-config"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  namespace              = "flink"
  session_cluster_name   = "flink-session-cluster"
  kafka_cluster_name     = "bookstore-kafka"
  kafka_broker_namespace = "kafka"

  sql_dir             = "${get_repo_root()}/deploy/kubernetes/flink/catalog-indexer"
  risk_scorer_sql_dir = "${get_repo_root()}/deploy/kubernetes/flink/risk-scorer"

  opensearch_secret_name = "catalog-indexer-opensearch"

  s3_endpoint          = "http://seaweedfs-s3.seaweedfs.svc.cluster.local:8333"
  checkpoint_s3_bucket = "flink-checkpoints"

  flink_jobs = {}
}

terraform {
  source = "${get_repo_root()}/infra//resource/flink"
}
