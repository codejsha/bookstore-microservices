variable "namespace" {
  description = "Flink namespace"
  type        = string
}

variable "session_cluster_name" {
  description = "FlinkDeployment (session cluster) name"
  type        = string
}

variable "kafka_cluster_name" {
  description = "Kafka cluster name"
  type        = string
}

variable "kafka_broker_namespace" {
  description = "Kafka broker namespace"
  type        = string
}

variable "flink_jobs" {
  description = "Flink session job definitions"
  type = map(object({
    jar_uri     = string
    parallelism = optional(number, 1)
    args        = optional(list(string), [])
  }))
}

variable "indexer_image" {
  description = "Custom Flink image for the catalog-indexer."
  type        = string
  default     = "harbor.example.com/bookstore/catalog-indexer:1.2.0"
}

variable "opensearch_secret_name" {
  description = "K8s Secret (same ns) with OpenSearch admin credentials for the sink"
  type        = string
}

variable "sql_dir" {
  description = "Absolute path to the catalog-indexer SQL scripts (deploy/kubernetes/flink/catalog-indexer)."
  type        = string
}

variable "s3_endpoint" {
  description = "SeaweedFS S3 endpoint used as the catalog-indexer durable checkpoint/savepoint store."
  type        = string
}

variable "checkpoint_s3_bucket" {
  description = "S3 bucket holding catalog-indexer Flink checkpoints/savepoints."
  type        = string
}

variable "vault_url" {
  description = "Vault URL"
  type        = string
}

variable "vault_auth_role" {
  type = string
}

variable "vault_k8s_jwt" {
  description = "Short-lived Kubernetes ServiceAccount token for Vault kubernetes-auth login"
  type        = string
  sensitive   = true
}

variable "risk_scorer_sql_dir" {
  description = "Directory holding the risk-scorer SQL scripts"
  type        = string
}

