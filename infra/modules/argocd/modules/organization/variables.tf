variable "argocd_address" {
  description = "ArgoCD address"
  type        = string
}

variable "admin_username" {
  description = "Admin username"
  type        = string
  sensitive   = true
}

variable "admin_password" {
  description = "Admin password"
  type        = string
  sensitive   = true
}

variable "gitea_fqdn" {
  description = "Gitea FQDN"
  type        = string
}
