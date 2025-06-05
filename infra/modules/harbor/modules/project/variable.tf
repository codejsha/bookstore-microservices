variable "project_name" {
  description = "Name of the Harbor project"
  type        = string
}

variable "project_users" {
  description = "Usernames to be added to the Harbor project"
  type = list(string)
}

variable "user_role" {
  description = "Role to assign to users"
  type        = string
}
