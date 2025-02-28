variable "argocd_address" {
  description = "ArgoCD address"
  type        = string
}

variable "admin_username" {
  description = "Username"
  type        = string
}

variable "admin_password" {
  description = "Password"
  type        = string
  sensitive   = true
}

variable "gitea_fqdn" {
  description = "Gitea FQDN"
  type        = string
}
