terraform {
  required_providers {
    harbor = {
      source = "goharbor/harbor"
    }
    vault = {
      source = "hashicorp/vault"
    }
    random = {
      source = "hashicorp/random"
    }
  }
}

resource "harbor_robot_account" "ci" {
  name        = var.robot_name
  description = "Tekton CI: pushes service images and helm charts"
  level       = "system"
  secret      = random_password.robot.result

  dynamic "permissions" {
    for_each = var.projects
    content {
      kind      = "project"
      namespace = permissions.value
      access {
        action   = "push"
        resource = "repository"
      }
      access {
        action   = "pull"
        resource = "repository"
      }
    }
  }
}

resource "random_password" "robot" {
  length  = 32
  special = false
}

resource "vault_kv_secret_v2" "registry" {
  mount = "kv"
  name  = "harbor/ci/credentials"
  data_json = jsonencode({
    username = harbor_robot_account.ci.full_name
    password = random_password.robot.result
  })
}
