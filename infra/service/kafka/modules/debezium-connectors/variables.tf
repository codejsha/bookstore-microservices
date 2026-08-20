variable "namespace" {
  description = "Namespace (same as Kafka broker + KafkaConnect)"
  type        = string
}

variable "kafka_cluster_name" {
  type = string
}

variable "connect_cluster_name" {
  description = "KafkaConnect cluster name that these connectors bind to"
  type        = string
}

variable "catalog_postgres_namespace" {
  description = "Namespace where catalog Postgres (CloudNativePG) lives"
  type        = string
}

variable "catalog_postgres_service" {
  description = "Catalog CNPG read-write service name"
  type        = string
}

variable "catalog_postgres_database" {
  description = "Catalog database name"
  type        = string
}

variable "debezium_secret_namespace" {
  description = "Namespace of the Debezium credential K8s Secret"
  type        = string
}

variable "debezium_secret_name" {
  description = "K8s Secret holding Debezium user username+password keys"
  type        = string
}
