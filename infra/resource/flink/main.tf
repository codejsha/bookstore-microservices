terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
    vault = {
      source  = "hashicorp/vault"
      version = "~> 5.10"
    }
  }
}

locals {
  kafka_bootstrap_server = "${var.kafka_cluster_name}-kafka-bootstrap.${var.kafka_broker_namespace}.svc:9093"
}

data "vault_kv_secret_v2" "opensearch_admin" {
  mount = "kv"
  name  = "opensearch/admin/credentials"
}

data "vault_kv_secret_v2" "harbor_pull" {
  mount = "kv"
  name  = "harbor/users/harbor-devops/credentials"
}

data "vault_kv_secret_v2" "seaweedfs_s3" {
  mount = "kv"
  name  = "seaweedfs/s3/credentials"
}

module "job" {
  for_each             = var.flink_jobs
  source               = "./modules/job"
  namespace            = var.namespace
  job_name             = each.key
  session_cluster_name = var.session_cluster_name
  job_jar_uri          = each.value.jar_uri
  job_parallelism      = each.value.parallelism
  job_args = concat(
    ["--bootstrap-servers", local.kafka_bootstrap_server],
    each.value.args
  )
}

module "catalog_indexer" {
  source                 = "./modules/catalog-indexer"
  namespace              = var.namespace
  indexer_image          = var.indexer_image
  opensearch_secret_name = var.opensearch_secret_name
  opensearch_username    = data.vault_kv_secret_v2.opensearch_admin.data["username"]
  opensearch_password    = data.vault_kv_secret_v2.opensearch_admin.data["password"]
  harbor_pull_username   = data.vault_kv_secret_v2.harbor_pull.data["username"]
  harbor_pull_password   = data.vault_kv_secret_v2.harbor_pull.data["password"]
  sql_dir                = var.sql_dir
  s3_endpoint            = var.s3_endpoint
  s3_bucket              = var.checkpoint_s3_bucket
  s3_access_key_id       = data.vault_kv_secret_v2.seaweedfs_s3.data["access_key_id"]
  s3_secret_key          = data.vault_kv_secret_v2.seaweedfs_s3.data["secret_access_key"]
}

module "risk_scorer" {
  source                 = "./modules/risk-scorer"
  namespace              = var.namespace
  runner_image           = var.indexer_image
  flink_service_account  = "flink-operator"
  image_pull_secret_name = "harbor-pull"
  sql_dir                = var.risk_scorer_sql_dir
  s3_endpoint            = var.s3_endpoint
  s3_bucket              = var.checkpoint_s3_bucket

  depends_on = [module.catalog_indexer]
}

