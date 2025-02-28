terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
  }
}

locals {
  namespace         = "bookstore"
  argocd_namespace  = "argocd"
  project           = local.namespace
  example_corp      = "example-corp"
  admin_service     = "admin"
  catalog_service   = "catalog"
  customer_service  = "customer"
  identity_service  = "identity"
  inventory_service = "inventory"
  order_service     = "order"
  payment_service   = "payment"
  git_ssh_url       = "gitea@${var.gitea_fqdn}"
}

resource "local_file" "ssh_private_key" {
  for_each = {
    for k, v in {
      (local.admin_service)     = data.vault_generic_secret.admin_helm_ssh_keys.data
      (local.catalog_service)   = data.vault_generic_secret.catalog_helm_ssh_keys.data
      (local.customer_service)  = data.vault_generic_secret.customer_helm_ssh_keys.data
      (local.identity_service)  = data.vault_generic_secret.identity_helm_ssh_keys.data
      (local.inventory_service) = data.vault_generic_secret.inventory_helm_ssh_keys.data
      (local.order_service)     = data.vault_generic_secret.order_helm_ssh_keys.data
      (local.payment_service)   = data.vault_generic_secret.payment_helm_ssh_keys.data
    } : k => v

  }
  content  = each.value["private"]
  filename = "${path.module}/keys/${each.key}_helm_repo.key"
}

resource "null_resource" "add_argocd_repositories" {
  depends_on = [local_file.ssh_private_key]
  provisioner "local-exec" {
    command = <<EOT
argocd login --username ${var.admin_username} --password ${var.admin_password} --insecure --plaintext ${var.argocd_address};
argocd repo add --upsert --name ${local.admin_service}-helm "${local.git_ssh_url}:${local.example_corp}/${local.admin_service}-helm.git" --insecure-ignore-host-key --ssh-private-key-path "${path.module}/keys/${local.admin_service}_helm_repo.key" --grpc-web;
argocd repo add --upsert --name ${local.catalog_service}-helm "${local.git_ssh_url}:${local.example_corp}/${local.catalog_service}-helm.git" --insecure-ignore-host-key --ssh-private-key-path "${path.module}/keys/${local.catalog_service}_helm_repo.key" --grpc-web;
argocd repo add --upsert --name ${local.customer_service}-helm "${local.git_ssh_url}:${local.example_corp}/${local.customer_service}-helm.git" --insecure-ignore-host-key --ssh-private-key-path "${path.module}/keys/${local.customer_service}_helm_repo.key" --grpc-web;
argocd repo add --upsert --name ${local.identity_service}-helm "${local.git_ssh_url}:${local.example_corp}/${local.identity_service}-helm.git" --insecure-ignore-host-key --ssh-private-key-path "${path.module}/keys/${local.identity_service}_helm_repo.key" --grpc-web;
argocd repo add --upsert --name ${local.inventory_service}-helm "${local.git_ssh_url}:${local.example_corp}/${local.inventory_service}-helm.git" --insecure-ignore-host-key --ssh-private-key-path "${path.module}/keys/${local.inventory_service}_helm_repo.key" --grpc-web;
argocd repo add --upsert --name ${local.order_service}-helm "${local.git_ssh_url}:${local.example_corp}/${local.order_service}-helm.git" --insecure-ignore-host-key --ssh-private-key-path "${path.module}/keys/${local.order_service}_helm_repo.key" --grpc-web;
argocd repo add --upsert --name ${local.payment_service}-helm "${local.git_ssh_url}:${local.example_corp}/${local.payment_service}-helm.git" --insecure-ignore-host-key --ssh-private-key-path "${path.module}/keys/${local.payment_service}_helm_repo.key" --grpc-web;
    EOT
  }
}

resource "kubernetes_manifest" "argocd_project" {
  depends_on = [null_resource.add_argocd_repositories]

  manifest = {
    apiVersion = "argoproj.io/v1alpha1"
    kind       = "AppProject"
    metadata = {
      name      = local.project
      namespace = local.argocd_namespace
      finalizers = [
        "resources-finalizer.argocd.argoproj.io"
      ]
    }
    spec = {
      destinations = [
        {
          namespace = local.namespace
          name      = "in-cluster"
          server    = "https://kubernetes.default.svc"
        }
      ]
      sourceRepos = [
        "${local.git_ssh_url}:${local.example_corp}/${local.admin_service}-helm.git",
        "${local.git_ssh_url}:${local.example_corp}/${local.catalog_service}-helm.git",
        "${local.git_ssh_url}:${local.example_corp}/${local.customer_service}-helm.git",
        "${local.git_ssh_url}:${local.example_corp}/${local.identity_service}-helm.git",
        "${local.git_ssh_url}:${local.example_corp}/${local.inventory_service}-helm.git",
        "${local.git_ssh_url}:${local.example_corp}/${local.order_service}-helm.git",
        "${local.git_ssh_url}:${local.example_corp}/${local.payment_service}-helm.git"
      ]
      clusterResourceWhitelist = [
        {
          group = "*"
          kind  = "*"
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "argocd_command_applications" {
  depends_on = [null_resource.add_argocd_repositories]

  for_each = {
    for app in [
      local.admin_service, local.catalog_service, local.customer_service, local.inventory_service,
      local.identity_service, local.order_service, local.payment_service
    ] : app => app
  }

  manifest = {
    apiVersion = "argoproj.io/v1alpha1"
    kind       = "Application"
    metadata = {
      name      = "${each.key}-command"
      namespace = local.argocd_namespace
    }
    spec = {
      project = local.project
      destination = {
        server    = "https://kubernetes.default.svc"
        namespace = local.namespace
      }
      source = {
        repoURL        = "${local.git_ssh_url}:${local.example_corp}/${each.key}-helm.git"
        path           = "."
        targetRevision = "HEAD"
        helm = {
          valueFiles = [
            "values.yaml",
            "values-command.yaml"
          ]
        }
      }
    }
  }
}

resource "kubernetes_manifest" "argocd_query_applications" {
  depends_on = [null_resource.add_argocd_repositories]
  for_each = {
    for app in [
      local.admin_service, local.catalog_service, local.customer_service, local.inventory_service,
      local.identity_service, local.order_service, local.payment_service
    ] : app => app
  }

  manifest = {
    apiVersion = "argoproj.io/v1alpha1"
    kind       = "Application"
    metadata = {
      name      = "${each.key}-query"
      namespace = local.argocd_namespace
    }
    spec = {
      project = kubernetes_manifest.argocd_project.manifest.metadata.name
      destination = {
        server    = "https://kubernetes.default.svc"
        namespace = local.namespace
      }
      source = {
        repoURL        = "${local.git_ssh_url}:${local.example_corp}/${each.key}-helm.git"
        path           = "."
        targetRevision = "HEAD"
        helm = {
          valueFiles = [
            "values.yaml",
            "values-query.yaml"
          ]
        }
      }
    }
  }
}
