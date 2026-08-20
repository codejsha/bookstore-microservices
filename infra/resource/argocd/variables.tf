variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "argocd_address" {
  description = "ArgoCD address"
  type        = string
}

variable "gitea_ssh_fqdn" {
  description = "Gitea SSH service FQDN used for the Argo CD repository connections"
  type        = string
}

variable "org_name" {
  description = "Organization name"
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

variable "app_repos" {
  description = "List of repositories"
  type        = list(string)
}

variable "environment" {
  description = "Deployment environment"
  type        = string
  default     = "dev"

  validation {
    condition     = contains(["dev", "prod"], var.environment)
    error_message = "environment must be one of: dev, prod."
  }
}

