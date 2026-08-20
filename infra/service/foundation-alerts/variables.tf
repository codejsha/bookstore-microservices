variable "vault_url" {
  description = "Vault URL"
  type        = string
}

variable "vault_k8s_jwt" {
  description = "Vault kubernetes-auth JWT (SA token used by the vault provider)"
  type        = string
  sensitive   = true
}

variable "vault_auth_role" {
  type = string
}

variable "grafana_url" {
  description = "Grafana base URL (external gateway host) for the alert-rule provider"
  type        = string
}
