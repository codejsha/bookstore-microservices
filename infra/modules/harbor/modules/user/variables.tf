variable "harbor_users" {
  description = "List of user mappings"
  type = set(object({
    user_name     = string
    user_password = string
  }))
}
