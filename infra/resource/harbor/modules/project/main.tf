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

resource "harbor_project_member_group" "bookstore_project_groups" {
  for_each = {
    for group in var.harbor_projects[local.bookstore_proj].groups : group.group_name => group
  }
  project_id = harbor_project.bookstore_project.id
  group_name = each.value.group_name
  type       = "oidc"
  role       = each.value.role
}

resource "harbor_project_member_group" "bookstore_helm_charts_project_groups" {
  for_each = {
    for group in var.harbor_projects[local.bookstore_helm_charts_proj].groups : group.group_name => group
  }
  project_id = harbor_project.bookstore_helm_charts_project.id
  group_name = each.value.group_name
  type       = "oidc"
  role       = each.value.role
}
