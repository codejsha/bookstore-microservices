variable "vault_url" {
  description = "Vault URL"
  type        = string
}

variable "local_output_dir" {
  description = "Absolute directory where the root and intermediate CA certificates are written"
  type        = string
}
