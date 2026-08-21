variable "harbor_projects" {
  description = "Harbor project configuration with members"
  type = map(object({
    project_name = string
    is_public    = bool
    members = list(object({
      username = string
      role     = string
    }))
  }))
}
