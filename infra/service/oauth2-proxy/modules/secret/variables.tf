variable "namespace" {
  description = "Namespace oauth2-proxy + the synced Secret live in"
  type        = string
}

variable "secret_name" {
  description = "Name of the destination Kubernetes Secret VSO creates"
  type        = string
}

variable "vault_kv_path" {
  description = "Vault KV v2 path (relative to the kv mount) to sync"
  type        = string
}
