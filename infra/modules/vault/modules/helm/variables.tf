variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "vault_tls_key" {
  description = "Vault TLS key"
  type        = string
}

variable "vault_tls_crt" {
  description = "Vault TLS certificate"
  type        = string
}

variable "vault_kube_ca_crt" {
  description = "Kubernetes CA certificate"
  type        = string
}
