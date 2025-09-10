variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "postgresql_username" {
  description = "PostgreSQL username for the gitea database"
  type        = string
}

variable "admin_username" {
  description = "Gitea admin username"
  type        = string
}

variable "admin_email" {
  description = "Gitea admin email"
  type        = string
}

variable "ssh_host_key_secret_name" {
  description = "Name of the Kubernetes secret holding Gitea's SSH host key, mounted into the Gitea container"
  type        = string
  default     = "gitea-ssh-host-key"
}
