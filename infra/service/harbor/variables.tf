variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "harbor_address" {
  description = "Harbor address"
  type        = string
}

variable "harbor_service_name" {
  description = "Harbor in-cluster Service name"
  type        = string
}

variable "vault_url" {
  description = "Vault URL"
  type        = string
}

variable "vault_auth_role" {
  type = string
}

variable "grafana_url" {
  description = "Grafana base URL (external gateway host) for the alert-rule provider"
  type        = string
}

variable "vault_k8s_jwt" {
  description = "Short-lived Kubernetes ServiceAccount token for Vault kubernetes-auth login"
  type        = string
  sensitive   = true
}

variable "aws_s3_api_url" {
  description = "AWS S3 API URL"
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
  type        = list(string)
}

variable "s3_registry_endpoint" {
  description = "In-cluster S3 endpoint the registry pods use for image/chart storage"
  type        = string
  default     = "http://seaweedfs-s3.seaweedfs.svc.cluster.local:8333"
}
