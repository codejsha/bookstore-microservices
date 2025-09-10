variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "argocd_address" {
  description = "ArgoCD address"
  type        = string
}

variable "argocd_service_name" {
  description = "ArgoCD in-cluster Service name"
  type        = string
}

variable "grafana_url" {
  description = "Grafana base URL (external gateway host) for the alert-rule provider"
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

variable "gitea_ssh_fqdn" {
  description = "Gitea SSH service FQDN the known_hosts entries are keyed on"
  type        = string
}
