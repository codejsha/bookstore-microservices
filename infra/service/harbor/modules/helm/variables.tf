variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "aws_access_key" {
  description = "AWS access key"
  type        = string
  sensitive   = true
}

variable "aws_secret_key" {
  description = "AWS secret key"
  type        = string
  sensitive   = true
}

variable "s3_bucket" {
  description = "S3 bucket for registry image/chart storage"
  type        = string
}

variable "s3_endpoint" {
  description = "In-cluster S3 endpoint for registry image/chart storage"
  type        = string
}
