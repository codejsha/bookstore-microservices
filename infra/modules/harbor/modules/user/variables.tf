variable "harbor_users" {
  description = "List of user mappings"
  type = set(object({
    username = string
    password = string
  }))
}
