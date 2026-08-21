variable "org_name" {
  description = "Organization name"
  type        = string
}

variable "team_name" {
  description = "Team name"
  type        = string
}

variable "user_credentials" {
  description = "List of user credentials"
  type = list(object({
    username = string
    password = string
  }))
}

variable "user_repos" {
  description = "List of repositories"
  type = list(string)
}
