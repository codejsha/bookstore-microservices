terraform {
  required_providers {
    argocd = {
      source = "argoproj-labs/argocd"
    }
  }
}

locals {
  project     = "bookstore"
  git_ssh_url = "git@${var.gitea_fqdn}"
  repos       = [
    for repo in var.app_repos :
    "${local.git_ssh_url}:${var.org_name}/${repo}.git"
  ]
}

resource "argocd_project" "bookstore" {
  metadata {
    name      = local.project
    namespace = var.namespace
  }
  spec {
    description  = "Bookstore project"
    source_repos = local.repos
    destination {
      server    = "https://kubernetes.default.svc"
      name      = "in-cluster"
      namespace = local.project
    }
    cluster_resource_whitelist {
      kind  = "*"
      group = "*"
    }
  }
}

resource "argocd_application" "bookstore_applications" {
  depends_on = [argocd_project.bookstore]
  for_each = toset(var.app_repos)
  metadata {
    name = replace(each.key, "-helm", "")
    namespace = var.namespace
  }
  spec {
    project = argocd_project.bookstore.metadata[0].name
    source {
      repo_url        = "${local.git_ssh_url}:${var.org_name}/${each.key}.git"
      target_revision = "HEAD"
      path            = "."
      helm {
        value_files = [
          "values.yaml",
          "values-dev.yaml"
        ]
      }
    }
    destination {
      server    = "https://kubernetes.default.svc"
      namespace = local.project
    }
  }
}
