variable "operator_namespace" {
  description = "Namespace name"
  type        = string
}

variable "broker_namespace" {
  description = "Namespace name"
  type        = string
}

variable "kafka_cluster_name" {
  description = "Kafka cluster name"
  type        = string
}

variable "connect_cluster_name" {
  description = "KafkaConnect cluster name"
  type        = string
}

variable "connect_image" {
  description = "Fully-qualified image ref where Strimzi pushes the built KafkaConnect image"
  type        = string
}

variable "harbor_registry_host" {
  description = "Harbor registry host (auth key in the dockerconfigjson pushSecret; must match connect_image's host)"
  type        = string
}

variable "replicas" {
  description = "KafkaConnect worker replicas"
  type        = number
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

variable "vault_url" {
  description = "Vault address for the kubernetes-auth provider login"
  type        = string
}

variable "vault_auth_role" {
  type = string
}

variable "vault_k8s_jwt" {
  description = "SA token for Vault kubernetes-auth (kubectl create token tf-kafka)"
  type        = string
  sensitive   = true
}

variable "grafana_url" {
  description = "Grafana base URL (external gateway host) for the alert-rule provider"
  type        = string
}
