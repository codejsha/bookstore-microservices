variable "repos" {
  description = "GitHub repo names (under github_owner) mirrored into Gitea under the same name."
  type        = set(string)
}

variable "github_owner" {
  description = "GitHub owner/org that holds the upstream library repos."
  type        = string
}

variable "github_branch" {
  description = "Upstream branch used as the change signal and as the mirrored repo's default branch."
  type        = string
  default     = "main"
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
