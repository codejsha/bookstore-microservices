variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "grafana_address" {
  description = "Grafana address"
  type        = string
}

variable "grafana_service_name" {
  description = "Grafana in-cluster Service name"
  type        = string
}

variable "alloy_grpc_address" {
  description = "Alloy OTLP gRPC address"
  type        = string
}

variable "alloy_http_address" {
  description = "Alloy OTLP HTTP address"
  type        = string
}

variable "alloy_gateway_service_name" {
  description = "Alloy gateway in-cluster Service name (external OTLP ingress target)"
  type        = string
}

variable "telemetry_route_namespaces" {
  description = "Namespaces whose HTTPRoutes may reference the alloy-gateway Service cross-namespace (ReferenceGrant). Used by SPA charts that proxy same-origin /faro/collect and /otlp paths to Alloy."
  type        = list(string)
  default     = ["bookstore"]
}

variable "s3_endpoint" {
  description = "S3 endpoint (SeaweedFS)"
  type        = string
}

variable "s3_region" {
  description = "S3 region"
  type        = string
}

variable "loki_s3_bucket" {
  description = "S3 bucket for Loki chunks/index"
  type        = string
}

variable "tempo_s3_bucket" {
  description = "S3 bucket for Tempo blocks"
  type        = string
}

variable "pyroscope_s3_bucket" {
  description = "S3 bucket for Pyroscope profiles"
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
