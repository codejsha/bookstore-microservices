variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "gitea_address" {
  description = "Gitea address"
  type        = string
}

variable "gitea_service_name" {
  description = "Gitea in-cluster Service name"
  type        = string
}

variable "admin_username" {
  description = "Admin username"
  type        = string
}

variable "admin_email" {
  description = "Admin email"
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

variable "grafana_url" {
  description = "Grafana base URL (external gateway host) for the alert-rule provider"
  type        = string
}
