variable "services" {
  description = "Service names under services/<name> to mirror"
  type        = list(string)
}

variable "repo_root" {
  description = "Absolute repo root; `git archive` runs against this working tree"
  type        = string
}

variable "repo_suffix" {
  description = "Suffix appended to the service name for its Gitea repo"
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
