variable "namespace" {
  description = "Namespace the Flink deployment runs in"
  type        = string
}

variable "runner_image" {
  description = "Flink SQL runner image (shared with catalog-indexer)"
  type        = string
}

variable "flink_service_account" {
  description = "ServiceAccount for the FlinkDeployment"
  type        = string
}

variable "image_pull_secret_name" {
  description = "Existing Harbor pull secret in the namespace (created by the catalog-indexer module)"
  type        = string
}

variable "s3_secret_name" {
  description = "Existing S3 credentials secret in the namespace (created by the catalog-indexer module)"
  type        = string
  default     = "catalog-indexer-s3"
}

variable "sql_dir" {
  description = "Directory holding the risk-scorer SQL scripts"
  type        = string
}

variable "s3_endpoint" {
  type = string
}

variable "s3_bucket" {
  type = string
}

variable "kafka_namespace" {
  description = "Namespace of the Strimzi Kafka cluster (KafkaTopic CR lives there)"
  type        = string
  default     = "kafka"
}

variable "kafka_cluster_name" {
  description = "Strimzi cluster name for the KafkaTopic label"
  type        = string
  default     = "bookstore-kafka"
}

variable "valkey_host" {
  description = "Identity's Valkey host the auto risk flags are written to"
  type        = string
  default     = "identity-valkey-headless.bookstore.svc.cluster.local"
}

variable "risk_score_threshold" {
  description = "5-minute windowed score at or above which a principal is auto-restricted"
  type        = number
  default     = 3000
}

variable "auto_flag_ttl_seconds" {
  description = "TTL of an auto-restrict flag"
  type        = number
  default     = 900
}
