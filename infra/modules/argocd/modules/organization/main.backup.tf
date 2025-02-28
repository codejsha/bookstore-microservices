# terraform {
#     required_providers {
#         argocd = {
#             source = "argoproj-labs/argocd"
#         }
#     }
# }
#
# locals {
#     namespace         = "bookstore"
#     argocd_namespace  = "argocd"
#     project           = local.namespace
#     example_corp      = "example-corp"
#     admin_service     = "admin"
#     catalog_service   = "catalog"
#     customer_service  = "customer"
#     inventory_service = "inventory"
#     order_service     = "order"
#     payment_service   = "payment"
#     git_ssh_url       = "gitea@${var.gitea_fqdn}"
# }
#
# resource "argocd_project" "bookstore" {
#     metadata {
#         name      = local.project
#         namespace = local.argocd_namespace
#     }
#     spec {
#         description = "Bookstore project"
#         source_repos = [
#             argocd_repository.admin_service_helm_repo.repo,
#             argocd_repository.catalog_service_helm_repo.repo,
#             argocd_repository.customer_service_helm_repo.repo,
#             argocd_repository.inventory_service_helm_repo.repo,
#             argocd_repository.order_service_helm_repo.repo,
#             argocd_repository.payment_service_helm_repo.repo
#         ]
#         destination {
#             server    = "https://kubernetes.default.svc"
#             name      = "in-cluster"
#             namespace = local.namespace
#         }
#         cluster_resource_whitelist {
#             kind  = "*"
#             group = "*"
#         }
#     }
# }
#
# resource "argocd_repository" "admin_service_helm_repo" {
#     name     = "${local.admin_service}-helm"
#     repo     = "${local.git_ssh_url}:${local.example_corp}/${local.admin_service}-helm.git"
#     username = "gitea"
#     ssh_private_key = file("${path.module}/keys/admin_service_helm_repo.key")
#     insecure = true
# }
#
# resource "argocd_repository" "catalog_service_helm_repo" {
#     name     = "${local.catalog_service}-helm"
#     repo     = "${local.git_ssh_url}:${local.example_corp}/${local.catalog_service}-helm.git"
#     username = "gitea"
#     ssh_private_key = file("${path.module}/keys/catalog_service_helm_repo.key")
#     insecure = true
# }
#
# resource "argocd_repository" "customer_service_helm_repo" {
#     name     = "${local.customer_service}-helm"
#     repo     = "${local.git_ssh_url}:${local.example_corp}/${local.customer_service}-helm.git"
#     username = "gitea"
#     ssh_private_key = file("${path.module}/keys/customer_service_helm_repo.key")
#     insecure = true
# }
#
# resource "argocd_repository" "inventory_service_helm_repo" {
#     name     = "${local.inventory_service}-helm"
#     repo     = "${local.git_ssh_url}:${local.example_corp}/${local.inventory_service}-helm.git"
#     username = "gitea"
#     ssh_private_key = file("${path.module}/keys/inventory_service_helm_repo.key")
#     insecure = true
# }
#
# resource "argocd_repository" "order_service_helm_repo" {
#     name     = "${local.order_service}-helm"
#     repo     = "${local.git_ssh_url}:${local.example_corp}/${local.order_service}-helm.git"
#     username = "gitea"
#     ssh_private_key = file("${path.module}/keys/order_service_helm_repo.key")
#     insecure = true
# }
#
# resource "argocd_repository" "payment_service_helm_repo" {
#     name     = "${local.payment_service}-helm"
#     repo     = "${local.git_ssh_url}:${local.example_corp}/${local.payment_service}-helm.git"
#     username = "gitea"
#     ssh_private_key = file("${path.module}/keys/payment_service_helm_repo.key")
#     insecure = true
# }
#
# resource "argocd_application" "admin_command_service" {
#     metadata {
#         name      = "${local.admin_service}-command"
#         namespace = local.argocd_namespace
#     }
#     spec {
#         project = argocd_project.bookstore.metadata.0.name
#         destination {
#             server    = "https://kubernetes.default.svc"
#             namespace = local.namespace
#         }
#         source {
#             repo_url        = argocd_repository.admin_service_helm_repo.repo
#             target_revision = "HEAD"
#             path            = "."
#             helm {
#                 value_files = [
#                     "${path.module}/values/values.yaml",
#                     "${path.module}/values/values-command.yaml"
#                 ]
#             }
#         }
#     }
# }
#
# resource "argocd_application" "admin_query_service" {
#     metadata {
#         name      = "${local.admin_service}-query"
#         namespace = local.argocd_namespace
#     }
#     spec {
#         project = argocd_project.bookstore.metadata.0.name
#         destination {
#             server    = "https://kubernetes.default.svc"
#             namespace = local.namespace
#         }
#         source {
#             repo_url        = argocd_repository.admin_service_helm_repo.repo
#             target_revision = "HEAD"
#             path            = "."
#             helm {
#                 value_files = [
#                     "${path.module}/values/values.yaml",
#                     "${path.module}/values/values-query.yaml"
#                 ]
#             }
#         }
#     }
# }
#
# resource "argocd_application" "catalog_command_service" {
#     metadata {
#         name      = "${local.catalog_service}-command"
#         namespace = local.argocd_namespace
#     }
#     spec {
#         project = argocd_project.bookstore.metadata.0.name
#         destination {
#             server    = "https://kubernetes.default.svc"
#             namespace = local.namespace
#         }
#         source {
#             repo_url        = argocd_repository.catalog_service_helm_repo.repo
#             target_revision = "HEAD"
#             path            = "."
#             helm {
#                 value_files = [
#                     "${path.module}/values/values.yaml",
#                     "${path.module}/values/values-command.yaml"
#                 ]
#             }
#         }
#     }
# }
#
# resource "argocd_application" "catalog_query_service" {
#     metadata {
#         name      = "${local.catalog_service}-query"
#         namespace = local.argocd_namespace
#     }
#     spec {
#         project = argocd_project.bookstore.metadata.0.name
#         destination {
#             server    = "https://kubernetes.default.svc"
#             namespace = local.namespace
#         }
#         source {
#             repo_url        = argocd_repository.catalog_service_helm_repo.repo
#             target_revision = "HEAD"
#             path            = "."
#             helm {
#                 value_files = [
#                     "${path.module}/values/values.yaml",
#                     "${path.module}/values/values-query.yaml"
#                 ]
#             }
#         }
#     }
# }
#
# resource "argocd_application" "customer_command_service" {
#     metadata {
#         name      = "${local.customer_service}-command"
#         namespace = local.argocd_namespace
#     }
#     spec {
#         project = argocd_project.bookstore.metadata.0.name
#         destination {
#             server    = "https://kubernetes.default.svc"
#             namespace = local.namespace
#         }
#         source {
#             repo_url        = argocd_repository.customer_service_helm_repo.repo
#             target_revision = "HEAD"
#             path            = "."
#             helm {
#                 value_files = [
#                     "${path.module}/values/values.yaml",
#                     "${path.module}/values/values-command.yaml"
#                 ]
#             }
#         }
#     }
# }
#
# resource "argocd_application" "customer_query_service" {
#     metadata {
#         name      = "${local.customer_service}-query"
#         namespace = local.argocd_namespace
#     }
#     spec {
#         project = argocd_project.bookstore.metadata.0.name
#         destination {
#             server    = "https://kubernetes.default.svc"
#             namespace = local.namespace
#         }
#         source {
#             repo_url        = argocd_repository.customer_service_helm_repo.repo
#             target_revision = "HEAD"
#             path            = "."
#             helm {
#                 value_files = [
#                     "${path.module}/values/values.yaml",
#                     "${path.module}/values/values-query.yaml"
#                 ]
#             }
#         }
#     }
# }
#
# resource "argocd_application" "inventory_command_service" {
#     metadata {
#         name      = "${local.inventory_service}-command"
#         namespace = local.argocd_namespace
#     }
#     spec {
#         project = argocd_project.bookstore.metadata.0.name
#         destination {
#             server    = "https://kubernetes.default.svc"
#             namespace = local.namespace
#         }
#         source {
#             repo_url        = argocd_repository.inventory_service_helm_repo.repo
#             target_revision = "HEAD"
#             path            = "."
#             helm {
#                 value_files = [
#                     "${path.module}/values/values.yaml",
#                     "${path.module}/values/values-command.yaml"
#                 ]
#             }
#         }
#     }
# }
#
# resource "argocd_application" "inventory_query_service" {
#     metadata {
#         name      = "${local.inventory_service}-query"
#         namespace = local.argocd_namespace
#     }
#     spec {
#         project = argocd_project.bookstore.metadata.0.name
#         destination {
#             server    = "https://kubernetes.default.svc"
#             namespace = local.namespace
#         }
#         source {
#             repo_url        = argocd_repository.inventory_service_helm_repo.repo
#             target_revision = "HEAD"
#             path            = "."
#             helm {
#                 value_files = [
#                     "${path.module}/values/values.yaml",
#                     "${path.module}/values/values-query.yaml"
#                 ]
#             }
#         }
#     }
# }
#
# resource "argocd_application" "order_command_service" {
#     metadata {
#         name      = "${local.order_service}-command"
#         namespace = local.argocd_namespace
#     }
#     spec {
#         project = argocd_project.bookstore.metadata.0.name
#         destination {
#             server    = "https://kubernetes.default.svc"
#             namespace = local.namespace
#         }
#         source {
#             repo_url        = argocd_repository.order_service_helm_repo.repo
#             target_revision = "HEAD"
#             path            = "."
#             helm {
#                 value_files = [
#                     "${path.module}/values/values.yaml",
#                     "${path.module}/values/values-command.yaml"
#                 ]
#             }
#         }
#     }
# }
#
# resource "argocd_application" "order_query_service" {
#     metadata {
#         name      = "${local.order_service}-query"
#         namespace = local.argocd_namespace
#     }
#     spec {
#         project = argocd_project.bookstore.metadata.0.name
#         destination {
#             server    = "https://kubernetes.default.svc"
#             namespace = local.namespace
#         }
#         source {
#             repo_url        = argocd_repository.order_service_helm_repo.repo
#             target_revision = "HEAD"
#             path            = "."
#             helm {
#                 value_files = [
#                     "${path.module}/values/values.yaml",
#                     "${path.module}/values/values-query.yaml"
#                 ]
#             }
#         }
#     }
# }
#
# resource "argocd_application" "payment_command_service" {
#     metadata {
#         name      = "${local.payment_service}-command"
#         namespace = local.argocd_namespace
#     }
#     spec {
#         project = argocd_project.bookstore.metadata.0.name
#         destination {
#             server    = "https://kubernetes.default.svc"
#             namespace = local.namespace
#         }
#         source {
#             repo_url        = argocd_repository.payment_service_helm_repo.repo
#             target_revision = "HEAD"
#             path            = "."
#             helm {
#                 value_files = [
#                     "${path.module}/values/values.yaml",
#                     "${path.module}/values/values-command.yaml"
#                 ]
#             }
#         }
#     }
# }
#
# resource "argocd_application" "payment_query_service" {
#     metadata {
#         name      = "${local.payment_service}-query"
#         namespace = local.argocd_namespace
#     }
#     spec {
#         project = argocd_project.bookstore.metadata.0.name
#         destination {
#             server    = "https://kubernetes.default.svc"
#             namespace = local.namespace
#         }
#         source {
#             repo_url        = argocd_repository.payment_service_helm_repo.repo
#             target_revision = "HEAD"
#             path            = "."
#             helm {
#                 value_files = [
#                     "${path.module}/values/values.yaml",
#                     "${path.module}/values/values-query.yaml"
#                 ]
#             }
#         }
#     }
# }
#
# resource "argocd_account_token" "devadmin_token" {
#     account    = "devadmin"
#     expires_in = "8760h"
# }
#
# resource "argocd_account_token" "devopsadmin_token" {
#     account    = "devopsadmin"
#     expires_in = "8760h"
# }
#
# resource "local_file" "devadmin_token" {
#     content  = argocd_account_token.devadmin_token.jwt
#     filename = "${path.module}/keys/devadmin_token.jwt"
# }
#
# resource "local_file" "devopsadmin_token" {
#     content  = argocd_account_token.devopsadmin_token.jwt
#     filename = "${path.module}/keys/devopsadmin_token.jwt"
# }
