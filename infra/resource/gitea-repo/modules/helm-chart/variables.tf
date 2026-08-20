variable "app_repos" {
  description = "List of <svc>-helm repositories to push charts for"
  type        = list(string)
}

variable "helm_dir" {
  description = "Absolute path to deploy/kubernetes/helm"
  type        = string
}

variable "org_name" {
  description = "Gitea organization name"
  type        = string
}

variable "gitea_api_url" {
  description = "Gitea API base URL"
  type        = string
}

variable "gitea_http_url" {
  description = "Gitea HTTP base URL for git push"
  type        = string
}

variable "gitea_username" {
  description = "Gitea admin username"
  type        = string
  sensitive   = true
}

variable "gitea_password" {
  description = "Gitea admin password/token"
  type        = string
  sensitive   = true
}
