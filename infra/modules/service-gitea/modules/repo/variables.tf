variable "gitea_username" {
  description = "Gitea username"
  type        = string
  sensitive   = true
}

variable "repo_name" {
  description = "Repository name"
  type        = string
}

variable "repo_description" {
  description = "Repository description"
  type        = string
}
