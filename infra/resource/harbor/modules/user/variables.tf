variable "harbor_users" {
  description = "List of user mappings"
  type = set(object({
    username = string
    email    = string
    password = string
  }))
  sensitive = true
}
