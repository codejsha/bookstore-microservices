variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "seaweedfs_access_key" {
  description = "SeaweedFS S3 access key"
  type        = string
  sensitive   = true
}

variable "seaweedfs_secret_key" {
  description = "SeaweedFS S3 secret key"
  type        = string
  sensitive   = true
}
