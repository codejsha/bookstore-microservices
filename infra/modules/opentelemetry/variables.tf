variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "otel_collector_grpc_address" {
  description = "OpenTelemetry Collector gRPC address"
}

variable "otel_collector_http_address" {
  description = "OpenTelemetry Collector HTTP address"
}

variable "otel_collector_fqdn" {
  description = "OpenTelemetry Collector FQDN"
}

variable "opensearch_username" {
  description = "OpenSearch username"
  type        = string
  sensitive   = true
}

variable "opensearch_password" {
  description = "OpenSearch password"
  type        = string
  sensitive   = true
}

variable "vault_url" {
  description = "Vault URL"
  type        = string
}

variable "vault_token" {
  description = "Vault authentication token"
  type        = string
  sensitive   = true
}

variable "kube_ca_crt_path" {
  description = "Kubernetes CA certificate file path"
  type        = string
}
