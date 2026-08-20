terraform {
  required_providers {
    argocd = {
      source = "argoproj-labs/argocd"
    }
  }
}

locals {
  project     = "bookstore"
  git_ssh_url = "ssh://git@${var.gitea_ssh_fqdn}:${var.gitea_ssh_port}/${var.org_name}"
  repos = [
    for repo in var.app_repos :
    "${local.git_ssh_url}/${repo}.git"
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
  for_each   = toset(var.app_repos)
  metadata {
    name      = replace(each.key, "-helm", "")
    namespace = var.namespace
  }
  spec {
    project = argocd_project.bookstore.metadata[0].name
    source {
      repo_url        = "${local.git_ssh_url}/${each.key}.git"
      target_revision = "HEAD"
      path            = "."
      helm {
        value_files = [
          "values.yaml",
          "values-${var.environment}.yaml"
        ]
      }
    }
    destination {
      server    = "https://kubernetes.default.svc"
      namespace = local.project
    }
    ignore_difference {
      group               = "gateway.networking.k8s.io"
      kind                = "HTTPRoute"
      jq_path_expressions = [".spec.rules[]?.backendRefs[]?.weight"]
    }

    ignore_difference {
      group               = ""
      kind                = "Service"
      jq_path_expressions = [".spec.selector[\"rollouts-pod-template-hash\"]"]
    }

    sync_policy {
      sync_options = [
        "CreateNamespace=false",
        "ServerSideApply=true",
      ]
    }
  }
}
