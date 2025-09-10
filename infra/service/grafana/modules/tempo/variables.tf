variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "s3_endpoint" {
  description = "S3 endpoint (SeaweedFS)"
  type        = string
}

variable "s3_region" {
  description = "S3 region"
  type        = string
}

variable "s3_bucket" {
  description = "S3 bucket for Tempo blocks"
  type        = string
}

variable "s3_access_key_id" {
  description = "S3 access key ID"
  type        = string
  sensitive   = true
}

variable "s3_secret_key" {
  description = "S3 secret access key"
  type        = string
  sensitive   = true
}
