variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "prometheus_address" {
  description = "Prometheus address"
  type        = string
}

variable "prometheus_service_name" {
  description = "Prometheus in-cluster Service name"
  type        = string
}

variable "vault_url" {
  description = "Vault URL"
  type        = string
}

variable "vault_auth_role" {
  type = string
}

variable "grafana_url" {
  description = "Grafana base URL (external gateway host) for the alert-rule provider"
  type        = string
}

variable "vault_k8s_jwt" {
  description = "Short-lived Kubernetes ServiceAccount token for Vault kubernetes-auth login"
  type        = string
  sensitive   = true
}
