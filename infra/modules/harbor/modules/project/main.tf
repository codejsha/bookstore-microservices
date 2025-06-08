terraform {
  required_providers {
    harbor = {
      source = "goharbor/harbor"
    }
  }
}

locals {
  bookstore_proj             = "bookstore"
  bookstore_helm_charts_proj = "bookstore-helm-charts"
}

resource "harbor_project" "bookstore_project" {
  name   = local.bookstore_proj
  public = var.harbor_projects[local.bookstore_proj].is_public
}

resource "harbor_project" "bookstore_helm_charts_project" {
  name   = local.bookstore_helm_charts_proj
  public = var.harbor_projects[local.bookstore_helm_charts_proj].is_public
}

resource "harbor_project_member_user" "bookstore_project_members" {
  for_each = {
    for user in var.harbor_projects[local.bookstore_proj].members : user.user_name => user
  }
  project_id = harbor_project.bookstore_project.id
  user_name  = each.value.user_name
  role       = each.value.user_role
}

resource "harbor_project_member_user" "bookstore_helm_charts_project_members" {
  for_each = {
    for user in var.harbor_projects[local.bookstore_helm_charts_proj].members : user.user_name => user
  }
  project_id = harbor_project.bookstore_helm_charts_project.id
  user_name  = each.value.user_name
  role       = each.value.user_role
}
