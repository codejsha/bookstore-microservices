include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "kubernetes" {
  path = find_in_parent_folders("_envcommon/provider_kubernetes.hcl")
}

include "helm" {
  path = find_in_parent_folders("_envcommon/provider_helm.hcl")
}

include "vault" {
  path = find_in_parent_folders("_envcommon/provider_vault.hcl")
}

dependencies {
  paths = [
    "../../service/prometheus-crds",
    "../../service/vault",
    "../../service/grafana",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role = "kafka"
}

inputs = {
  operator_namespace = "kafka-operator"
  broker_namespace   = "kafka"
  kafka_cluster_name = "bookstore-kafka"

  connect_image        = "harbor.example.com/bookstore/kafka-connect:latest"
  harbor_registry_host = "harbor.example.com"
  connect_cluster_name = "kafka-connect"
  replicas             = 1

  catalog_postgres_namespace = "bookstore"
  catalog_postgres_service   = "catalog-postgres-rw"
  catalog_postgres_database  = "catalog_db"
  debezium_secret_namespace  = "bookstore"
  debezium_secret_name       = "catalog-postgres-debezium"

  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  grafana_url = "https://grafana.example.com"
}

terraform {
  source = "${get_repo_root()}/infra//service/kafka"
}
