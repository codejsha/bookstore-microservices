terraform {
  required_providers {
    argocd = {
      source = "argoproj-labs/argocd"
    }
    vault = {
      source = "hashicorp/vault"
    }
  }
}

locals {
  project     = "bookstore"
  git_ssh_url = "gitea@${var.gitea_fqdn}"
  services = {
    catalog   = data.vault_generic_secret.catalog_helm_ssh_keys.data
    customer  = data.vault_generic_secret.customer_helm_ssh_keys.data
    identity  = data.vault_generic_secret.identity_helm_ssh_keys.data
    inventory = data.vault_generic_secret.inventory_helm_ssh_keys.data
    order     = data.vault_generic_secret.order_helm_ssh_keys.data
    payment   = data.vault_generic_secret.payment_helm_ssh_keys.data
  }
  repos = [
    for service in var.app_repos :
    "${local.git_ssh_url}:${var.org_name}/${service}-helm.git"
  ]
}

resource "argocd_repository" "service_helm_repos" {
  for_each = toset(var.app_repos)
  name            = "${each.key}-helm"
  repo            = "${local.git_ssh_url}:${var.org_name}/${each.key}-helm.git"
  username        = "gitea"
  ssh_private_key = local.services[each.key]["private"]
  insecure        = true
}

resource "argocd_project" "bookstore" {
  depends_on = [argocd_repository.service_helm_repos]
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

resource "argocd_application" "service_command_apps" {
  depends_on = [argocd_repository.service_helm_repos, argocd_project.bookstore]
  for_each = toset(var.app_repos)
  metadata {
    name      = "${each.key}-command"
    namespace = var.namespace
  }
  spec {
    project = argocd_project.bookstore.metadata[0].name
    source {
      repo_url        = "${local.git_ssh_url}:${var.org_name}/${each.key}-helm.git"
      target_revision = "HEAD"
      path            = "."
      helm {
        value_files = [
          "values.yaml",
          "values-command.yaml"
        ]
      }
    }
    destination {
      server    = "https://kubernetes.default.svc"
      namespace = local.project
    }
  }
}

resource "argocd_application" "service_query_apps" {
  depends_on = [argocd_repository.service_helm_repos, argocd_project.bookstore]
  for_each = toset(var.app_repos)
  metadata {
    name      = "${each.key}-query"
    namespace = var.namespace
  }
  spec {
    project = argocd_project.bookstore.metadata[0].name
    source {
      repo_url        = "${local.git_ssh_url}:${var.org_name}/${each.key}-helm.git"
      target_revision = "HEAD"
      path            = "."
      helm {
        value_files = [
          "values.yaml",
          "values-query.yaml"
        ]
      }
    }
    destination {
      server    = "https://kubernetes.default.svc"
      namespace = local.project
    }
  }
}
