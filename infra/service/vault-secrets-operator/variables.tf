variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "vault_address" {
  description = "In-cluster Vault address for the default VaultConnection (VSO runs in-cluster, so this is the ClusterIP service URL, not the external gateway host)."
  type        = string
}
