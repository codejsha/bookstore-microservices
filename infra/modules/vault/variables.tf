variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "vault_url" {
  description = "Vault URL"
  type        = string
}

variable "vault_address" {
  description = "Vault address"
  type        = string
}

variable "vault_fqdn" {
  description = "Vault FQDN"
  type        = string
}

variable "kube_ca_crt_path" {
  description = "Kubernetes CA certificate file path"
  type        = string
}

variable "kube_api_server_address" {
  description = "Kubernetes API server address"
  type        = string
}
