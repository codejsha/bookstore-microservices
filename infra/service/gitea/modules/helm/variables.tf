variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "admin_email" {
  description = "Admin email"
  type        = string
  sensitive   = true
}

variable "valkey_password" {
  description = "Valkey password"
  type        = string
  sensitive   = true
}

variable "postgresql_username" {
  description = "PostgreSQL username"
  type        = string
  sensitive   = true
}

variable "postgresql_password" {
  description = "PostgreSQL password"
  type        = string
  sensitive   = true
}

variable "ssh_host_key_secret_name" {
  description = "Kubernetes secret carrying the SSH host key mounted at /etc/gitea/ssh"
  type        = string
}
