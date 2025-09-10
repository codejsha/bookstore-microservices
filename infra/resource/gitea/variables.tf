variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "gitea_url" {
  description = "Gitea URL"
  type        = string
}

variable "admin_username" {
  description = "Admin username"
  type        = string
  sensitive   = true
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

variable "org_name" {
  description = "Organization name"
  type        = string
}

variable "dev_usernames" {
  description = "Development usernames; passwords are generated and stored in vault."
  type        = list(string)
}

variable "devops_usernames" {
  description = "DevOps usernames; passwords are generated and stored in vault."
  type        = list(string)
}

variable "dev_repos" {
  description = "List of repositories"
  type        = list(string)
}

variable "devops_repos" {
  description = "List of repositories"
  type        = list(string)
}
