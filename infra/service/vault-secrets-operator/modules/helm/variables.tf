variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "vault_address" {
  description = "In-cluster Vault address for the default VaultConnection."
  type        = string
}
