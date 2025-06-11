variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "initial_admin_password" {
  description = "Admin password"
  type        = string
  sensitive   = true
}
