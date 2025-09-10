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

variable "repo_name" {
  description = "Gitea repo name for the mirrored catalog"
  type        = string
}

variable "upstream_url" {
  description = "Upstream git URL to mirror"
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
