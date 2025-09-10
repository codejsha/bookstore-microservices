variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "org_name" {
  description = "Organization name"
  type        = string
}

variable "gitea_ssh_fqdn" {
  description = "Gitea SSH service FQDN"
  type        = string
}

variable "gitea_ssh_port" {
  description = "Gitea SSH service port"
  type        = number
  default     = 22
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

