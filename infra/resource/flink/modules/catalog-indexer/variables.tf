variable "namespace" {
  description = "Namespace where the catalog-indexer Flink job runs"
  type        = string
}

variable "indexer_image" {
  description = "Custom Flink image bundling the SQL-runner jar + Kafka/OpenSearch/Avro connectors."
  type        = string
  default     = "harbor.example.com/bookstore/catalog-indexer:1.0.0"
}

variable "flink_service_account" {
  description = "ServiceAccount the FlinkDeployment pods run as (must be RBAC-bound for the operator)."
  type        = string
  default     = "flink-operator"
}

variable "image_pull_secret_name" {
  description = "Name of the dockerconfigjson pull secret created for the private indexer image."
  type        = string
  default     = "harbor-pull"
}

variable "registry_host" {
  description = "Container registry host the pull secret authenticates to."
  type        = string
  default     = "harbor.example.com"
}

variable "harbor_pull_username" {
  description = "Harbor pull identity username."
  type        = string
}

variable "harbor_pull_password" {
  description = "Harbor pull identity password."
  type        = string
  sensitive   = true
}

variable "opensearch_secret_name" {
  description = "K8s Secret (same ns) with OpenSearch admin credentials for the sink"
  type        = string
}

variable "opensearch_username" {
  description = "OpenSearch admin username — materialised into opensearch_secret_name"
  type        = string
}

variable "opensearch_password" {
  description = "OpenSearch admin password — materialised into opensearch_secret_name"
  type        = string
  sensitive   = true
}

variable "sql_dir" {
  description = "Absolute path to the catalog-indexer SQL scripts directory."
  type        = string
}

variable "s3_endpoint" {
  description = "SeaweedFS S3 endpoint used as the durable checkpoint/savepoint store."
  type        = string
}

variable "s3_bucket" {
  description = "S3 bucket that holds catalog-indexer checkpoints/savepoints."
  type        = string
}

variable "s3_access_key_id" {
  description = "SeaweedFS S3 access key id."
  type        = string
  sensitive   = true
}

variable "s3_secret_key" {
  description = "SeaweedFS S3 secret access key."
  type        = string
  sensitive   = true
}
