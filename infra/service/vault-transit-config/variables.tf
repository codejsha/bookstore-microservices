variable "main_vault_namespace" {
  description = "Namespace of the MAIN Vault, where the auto-unseal token Secret is written"
  type        = string
  default     = "vault"
}

variable "vault_unsealer_url" {
  description = "Address of the unsealer Vault for the Terraform provider (laptop-reachable, e.g. port-forward to 127.0.0.1:8200)"
  type        = string
}

variable "vault_unsealer_token" {
  description = "Root/admin token of the unsealer Vault (transit provisioning only). Source from VAULT_UNSEALER_TOKEN env — never commit."
  type        = string
  sensitive   = true
}
