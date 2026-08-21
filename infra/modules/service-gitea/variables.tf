variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "gitea_url" {
  description = "Gitea URL"
  type        = string
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

variable "org_name" {
  description = "Organization name"
  type        = string
}

variable "dev_user_credentials" {
  description = "List of development user credentials"
  type = list(object({
    username = string
    password = string
  }))
}

variable "devops_user_credentials" {
  description = "List of DevOps user credentials"
  type = list(object({
    username = string
    password = string
  }))
}

variable "dev_repos" {
  description = "List of repositories"
  type = list(string)
}

variable "devops_repos" {
  description = "List of repositories"
  type = list(string)
}
