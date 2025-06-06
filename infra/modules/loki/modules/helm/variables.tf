variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "minio_username" {
  description = "Username"
  type        = string
}

variable "minio_password" {
  description = "Password"
  type        = string
  sensitive   = true
}
