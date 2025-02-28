variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "harbor_url" {
  description = "Harbor URL"
  type        = string
}

variable "harbor_address" {
  description = "Harbor address"
  type        = string
}

variable "harbor_fqdn" {
  description = "Harbor FQDN"
  type        = string
}

variable "harbor_username" {
  description = "Harbor username"
  type        = string
  sensitive   = true
}
variable "harbor_password" {
  description = "Harbor password"
  type        = string
  sensitive   = true
}

variable "name_prefix" {
  description = "Resource name prefix"
  type        = string
}

variable "minio_api_url" {
  description = "Minio API URL"
  type        = string
}

variable "aws_access_key" {
  description = "AWS Access Key"
  type        = string
  sensitive   = true
}

variable "aws_secret_key" {
  description = "AWS Secret Key"
  type        = string
  sensitive   = true
}

variable "bucket_names" {
  description = "List of bucket names"
  type = list(string)
}
