variable "webhooks" {
  description = "Map of Gitea source repo name => EventListener name. Keys drive which repos are swept for leftover el-* hooks."
  type        = map(string)
}

variable "org_name" {
  description = "Gitea organization that owns the source repos."
  type        = string
}

variable "source_repo_suffix" {
  description = "Suffix on the Gitea source repo name (must match the gitea-repo unit's source_repo_suffix)."
  type        = string
}

variable "gitea_api_url" {
  description = "Gitea API base URL."
  type        = string
}

variable "gitea_username" {
  description = "Gitea admin username."
  type        = string
  sensitive   = true
}

variable "gitea_password" {
  description = "Gitea admin password/token."
  type        = string
  sensitive   = true
}
