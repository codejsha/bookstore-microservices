variable "harbor_projects" {
  description = "Harbor project configuration with OIDC group members"
  type = map(object({
    project_name = string
    is_public    = bool
    groups = list(object({
      group_name = string
      role       = string
    }))
  }))
}
