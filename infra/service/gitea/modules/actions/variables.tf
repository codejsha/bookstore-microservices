variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "gitea_service_name" {
  description = "Gitea in-cluster Service name the runner registers against"
  type        = string
}

variable "gitea_api_url" {
  description = "Gitea API base URL reachable from the workstation (for minting the runner registration token)"
  type        = string
}

variable "admin_username" {
  description = "Gitea admin username"
  type        = string
  sensitive   = true
}

variable "admin_password" {
  description = "Gitea admin password"
  type        = string
  sensitive   = true
}

variable "token_secret_name" {
  description = "Name of the Kubernetes secret holding the runner registration token"
  type        = string
  default     = "act-runner-token"
}
