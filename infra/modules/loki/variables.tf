variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "minio_username" {
  description = "MinIO username"
  type        = string
  sensitive   = true
}

variable "minio_password" {
  description = "MinIO password"
  type        = string
  sensitive   = true
}
