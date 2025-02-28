variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "admin_username" {
  description = "Username"
  type        = string
}

variable "admin_password" {
  description = "Password"
  type        = string
  sensitive   = true
}
