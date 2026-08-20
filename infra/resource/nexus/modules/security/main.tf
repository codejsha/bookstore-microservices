terraform {
  required_providers {
    nexus = {
      source = "datadrivers/nexus"
    }
  }
}

locals {
  hosted_privileges = [
    "nx-repository-view-npm-${var.hosted_repository_name}-browse",
    "nx-repository-view-npm-${var.hosted_repository_name}-read",
    "nx-repository-view-npm-${var.hosted_repository_name}-add",
    "nx-repository-view-npm-${var.hosted_repository_name}-edit",
  ]

  group_privileges = [
    "nx-repository-view-npm-${var.group_repository_name}-browse",
    "nx-repository-view-npm-${var.group_repository_name}-read",
  ]

  maven_hosted_privileges = [
    "nx-repository-view-maven2-${var.maven_hosted_repository_name}-browse",
    "nx-repository-view-maven2-${var.maven_hosted_repository_name}-read",
    "nx-repository-view-maven2-${var.maven_hosted_repository_name}-add",
    "nx-repository-view-maven2-${var.maven_hosted_repository_name}-edit",
  ]

  maven_group_privileges = [
    "nx-repository-view-maven2-${var.maven_group_repository_name}-browse",
    "nx-repository-view-maven2-${var.maven_group_repository_name}-read",
  ]

  pypi_group_privileges = [
    "nx-repository-view-pypi-${var.pypi_group_repository_name}-browse",
    "nx-repository-view-pypi-${var.pypi_group_repository_name}-read",
  ]

  raw_hosted_privileges = [
    "nx-repository-view-raw-${var.raw_hosted_repository_name}-browse",
    "nx-repository-view-raw-${var.raw_hosted_repository_name}-read",
    "nx-repository-view-raw-${var.raw_hosted_repository_name}-add",
    "nx-repository-view-raw-${var.raw_hosted_repository_name}-edit",
  ]

  go_group_privileges = [
    "nx-repository-view-go-${var.go_group_repository_name}-browse",
    "nx-repository-view-go-${var.go_group_repository_name}-read",
  ]

  docker_proxy_privileges = [
    "nx-repository-view-docker-${var.docker_proxy_repository_name}-browse",
    "nx-repository-view-docker-${var.docker_proxy_repository_name}-read",
  ]
}

resource "nexus_security_anonymous" "anonymous" {
  enabled = false
}

resource "nexus_security_role" "publisher" {
  roleid      = var.role_id
  name        = var.role_id
  description = "Publish (view/add/edit) on the npm/maven/raw hosted repos, read on the npm/maven/pypi/go groups and the docker proxy"
  privileges = concat(
    local.hosted_privileges,
    local.group_privileges,
    local.maven_hosted_privileges,
    local.maven_group_privileges,
    local.pypi_group_privileges,
    local.raw_hosted_privileges,
    local.go_group_privileges,
    local.docker_proxy_privileges,
  )
}

resource "nexus_security_user" "publisher" {
  userid    = var.publisher_username
  firstname = "Bookstore"
  lastname  = "Publisher"
  email     = "${var.publisher_username}@example.com"
  status    = "active"
  roles     = [nexus_security_role.publisher.roleid]

  password_wo         = var.publisher_password
  password_wo_version = 1
}
