variable "org_name" {
  description = "Organization name"
  type        = string
}

variable "dev_users" {
  description = "List of users"
  type = list(string)
}

variable "dev_repos" {
  description = "List of repositories"
  type = list(string)
}

variable "devops_users" {
  description = "List of users"
  type = list(string)
}

variable "devops_repos" {
  description = "List of repositories"
  type = list(string)
}
