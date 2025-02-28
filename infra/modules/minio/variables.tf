variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "minio_api_address" {
  description = "MinIO API address"
  type        = string
}

variable "minio_ui_address" {
  description = "MinIO UI address"
  type        = string
}

variable "minio_fqdn" {
  description = "MinIO FQDN"
  type        = string
}

variable "minio_username" {
  description = "Username"
  type        = string
}

variable "minio_password" {
  description = "Password"
  type        = string
  sensitive   = true
}

variable "vault_url" {
  description = "Address"
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
