terraform {
  required_providers {
    harbor = {
      source = "goharbor/harbor"
    }
  }
}

resource "harbor_project" "project" {
  name   = var.project_name
  public = true
}

resource "harbor_project_member_user" "project_member" {
  for_each = toset(var.project_users)
  project_id = harbor_project.project.id
  role       = var.user_role
  user_name  = each.value
}
