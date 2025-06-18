variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "mysql_services" {
  description = "MySQL services"
  type = list(string)
}

variable "mysql_root_password" {
  description = "MySQL root password"
  type        = string
  sensitive   = true
}

variable "mysql_username" {
  description = "MySQL username"
  type        = string
  sensitive   = true
}

variable "mysql_password" {
  description = "MySQL password"
  type        = string
  sensitive   = true
}

variable "vault_url" {
  description = "Vault URL"
  type        = string
}

variable "vault_token" {
  description = "Vault authentication token"
  type        = string
  sensitive   = true
}

variable "kube_ca_crt_path" {
  description = "Kubernetes CA certificate file path"
  type        = string
}
