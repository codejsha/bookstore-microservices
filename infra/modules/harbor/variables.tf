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

variable "vault_url" {
  description = "Vault URL"
  type        = string
}

variable "vault_token" {
  description = "Vault authentication token"
  type        = string
  sensitive   = true
}

variable "kube_ca_crt_path" {
  description = "Kubernetes CA certificate file path"
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

variable "harbor_users" {
  description = "List of user mappings"
  type = set(object({
    user_name     = string
    user_password = string
  }))
}

variable "harbor_projects" {
  description = "Harbor project configuration with members"
  type = map(object({
    project_name = string
    is_public    = bool
    members = list(object({
      user_name = string
      user_role = string
    }))
  }))
}
