variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "admin_email" {
  description = "Operator email (recorded in Vault KV metadata)"
  type        = string
}
