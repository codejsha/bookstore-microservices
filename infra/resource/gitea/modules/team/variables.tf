variable "org_name" {
  description = "Organization name"
  type        = string
}

variable "team_name" {
  description = "Team name"
  type        = string
}


variable "user_repos" {
  description = "List of repositories"
  type        = list(string)
}

variable "include_all_repositories" {
  description = "Grant the team access to every current and future repository in the organization"
  type        = bool
  default     = false
}
