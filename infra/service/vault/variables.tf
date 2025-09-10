variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "vault_url" {
  description = "Vault URL (scheme + host[:port], used by the vault provider)"
  type        = string
}

variable "vault_address" {
  description = "Vault hostname (Gateway API HTTPRoute hostname + cert SAN)"
  type        = string
}

variable "vault_service_name" {
  description = "Vault in-cluster Service name (HTTPRoute backendRef)"
  type        = string
}

variable "kube_api_server_address" {
  description = "Kubernetes API server address (host:port, no scheme)"
  type        = string
}
