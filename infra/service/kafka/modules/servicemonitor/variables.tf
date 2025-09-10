variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "kafka_cluster_name" {
  description = "Kafka cluster name"
  type        = string
}

variable "operator_namespace" {
  description = "Namespace the Strimzi cluster operator runs in"
  type        = string
}
