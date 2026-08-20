variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "request_timeout" {
  description = "Per-request timeout the waypoint enforces on the hyperswitch server API. Same role as bookstore-mesh's request_timeout: the platform's only request deadline."
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
