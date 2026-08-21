variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "opensearch_address" {
  description = "OpenSearch address"
  type        = string
}

variable "opensearch_fqdn" {
  description = "OpenSearch FQDN"
  type        = string
}

variable "initial_admin_password" {
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
