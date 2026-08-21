variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "admin_email" {
  description = "Admin email"
  type        = string
  sensitive   = true
}

variable "admin_username" {
  description = "Admin username"
  type        = string
  sensitive   = true
}

variable "admin_password" {
  description = "Admin password"
  type        = string
  sensitive   = true
}

variable "valkey_password" {
  description = "Valkey password"
  type        = string
  sensitive   = true
}
