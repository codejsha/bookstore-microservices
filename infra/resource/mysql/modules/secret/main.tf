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
  creds_by_service = { for c in var.mysql_credentials : c.service => c }
}

resource "kubernetes_secret_v1" "mysql_cluster_secret" {
  for_each = toset(var.mysql_services)
  metadata {
    name      = "${each.key}-mysql-cluster-secret"
    namespace = var.namespace
    annotations = {
      "reflector.v1.k8s.emberstack.com/reflection-allowed"            = "true"
      "reflector.v1.k8s.emberstack.com/reflection-allowed-namespaces" = var.ci_namespace
      "reflector.v1.k8s.emberstack.com/reflection-auto-enabled"       = "true"
      "reflector.v1.k8s.emberstack.com/reflection-auto-namespaces"    = var.ci_namespace
    }
  }
  type = "Opaque"
  data = {
    rootUser     = "root"
    rootHost     = "%"
    rootPassword = var.mysql_root_password
  }
}

resource "vault_policy" "bookstore_mysql" {
  name   = "bookstore-mysql"
  policy = file("${path.module}/policy.hcl")
}

resource "vault_kubernetes_auth_backend_role" "bookstore_mysql" {
  role_name                        = "bookstore-mysql-role"
  bound_service_account_names      = var.mysql_service_accounts
  bound_service_account_namespaces = [var.namespace]
  token_policies                   = [vault_policy.bookstore_mysql.name]
  token_ttl                        = "3600"
}

resource "vault_kv_secret_v2" "bookstore_mysql" {
  for_each = toset(var.mysql_services)
  name     = "bookstore/${each.key}/mysql"
  mount    = "kv"
  data_json = jsonencode({
    root_password = var.mysql_root_password,
    username      = "${each.key}_app",
    password      = local.creds_by_service[each.key].password,
  })
}

variable "ci_namespace" {
  description = "Namespace where Tekton PipelineRuns execute; the DB credentials are reflected there for db-migrate"
  type        = string
  default     = "bookstore-ci"
}
