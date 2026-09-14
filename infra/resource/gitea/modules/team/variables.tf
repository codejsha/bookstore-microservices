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
