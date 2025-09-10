variable "repo_name" {
  description = "Gitea repo name for the custom Tekton catalog"
  type        = string
}

variable "content_dir" {
  description = "Absolute path to infra/resource/tekton/catalog"
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
