variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "gitea_address" {
  description = "Gitea address"
  type        = string
}

variable "gitea_fqdn" {
  description = "Gitea FQDN"
  type        = string
}

variable "admin_email" {
  description = "Admin email"
  type        = string
  sensitive   = true
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

variable "valkey_password" {
  description = "Valkey password"
  type        = string
  sensitive   = true
}
