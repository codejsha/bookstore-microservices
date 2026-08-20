terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

locals {
  creds_by_service = { for c in var.postgres_credentials : c.service => c }
}

resource "kubernetes_secret_v1" "postgres_app" {
  for_each = toset(var.postgres_services)
  metadata {
    name      = "${each.key}-postgres-app"
    namespace = var.namespace
    annotations = {
      "reflector.v1.k8s.emberstack.com/reflection-allowed"            = "true"
      "reflector.v1.k8s.emberstack.com/reflection-allowed-namespaces" = var.ci_namespace
      "reflector.v1.k8s.emberstack.com/reflection-auto-enabled"       = "true"
      "reflector.v1.k8s.emberstack.com/reflection-auto-namespaces"    = var.ci_namespace
    }
  }
  type = "kubernetes.io/basic-auth"
  data = {
    username = local.creds_by_service[each.key].username
    password = local.creds_by_service[each.key].password
  }
}

resource "kubernetes_secret_v1" "postgres_debezium" {
  for_each = toset(var.postgres_services)
  metadata {
    name      = "${each.key}-postgres-debezium"
    namespace = var.namespace
  }
  type = "kubernetes.io/basic-auth"
  data = {
    username = var.debezium_username
    password = var.debezium_password
  }
}

resource "vault_policy" "bookstore_postgres" {
  name   = "bookstore-postgres"
  policy = file("${path.module}/policy.hcl")
}

resource "vault_kubernetes_auth_backend_role" "bookstore_postgres" {
  role_name                        = "bookstore-postgres-role"
  bound_service_account_names      = var.postgres_service_accounts
  bound_service_account_namespaces = [var.namespace]
  token_policies                   = [vault_policy.bookstore_postgres.name]
  token_ttl                        = "3600"
}

resource "vault_policy" "bookstore_catalog" {
  name   = "bookstore-catalog"
  policy = file("${path.module}/catalog-policy.hcl")
}

resource "vault_kubernetes_auth_backend_role" "bookstore_catalog" {
  role_name                        = "bookstore-catalog-role"
  bound_service_account_names      = ["catalog-postgres"]
  bound_service_account_namespaces = [var.namespace]
  token_policies                   = [vault_policy.bookstore_catalog.name]
  token_ttl                        = "3600"
}

resource "vault_policy" "bookstore_identity" {
  name   = "bookstore-identity"
  policy = file("${path.module}/identity-policy.hcl")
}

resource "vault_kubernetes_auth_backend_role" "bookstore_identity" {
  role_name                        = "bookstore-identity-role"
  bound_service_account_names      = ["identity-postgres"]
  bound_service_account_namespaces = [var.namespace]
  token_policies                   = [vault_policy.bookstore_identity.name]
  token_ttl                        = "3600"
}

resource "vault_kv_secret_v2" "bookstore_postgres" {
  for_each = toset(var.postgres_services)
  name     = "bookstore/${each.key}/postgres"
  mount    = "kv"
  data_json = jsonencode({
    username = local.creds_by_service[each.key].username,
    password = local.creds_by_service[each.key].password,
  })
}

resource "vault_kv_secret_v2" "catalog_postgres_debezium" {
  name  = "bookstore/catalog/postgres-debezium"
  mount = "kv"
  data_json = jsonencode({
    username = var.debezium_username,
    password = var.debezium_password,
  })
}

variable "ci_namespace" {
  description = "Namespace where Tekton PipelineRuns execute; the DB credentials are reflected there for db-migrate"
  type        = string
  default     = "bookstore-ci"
}
