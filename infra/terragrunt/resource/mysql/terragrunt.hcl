include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "kubernetes" {
  path = find_in_parent_folders("_envcommon/provider_kubernetes.hcl")
}

include "vault" {
  path = find_in_parent_folders("_envcommon/provider_vault.hcl")
}

dependencies {
  paths = [
    "../bookstore-base",
    "../../service/mysql",
    "../../service/vault",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role = "mysql-config"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
    "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  namespace = "bookstore"
  mysql_services = ["order", "payment", "delivery", "notification", "support", "settlement"]

  mysql_client_image = "mysql:8.4"

  mysql_db_config = {
    order = {
      cluster_name  = "order-mysql"
      database_name = "order_db"
      secret_name   = "order-mysql-cluster-secret"
      app_user      = "order_app"
      instances     = 1
      version       = "8.4.3"
      storage_size  = "10Gi"
    }
    payment = {
      cluster_name  = "payment-mysql"
      database_name = "payment_db"
      secret_name   = "payment-mysql-cluster-secret"
      app_user      = "payment_app"
      instances     = 1
      version       = "8.4.3"
      storage_size  = "10Gi"
    }
    delivery = {
      cluster_name  = "delivery-mysql"
      database_name = "delivery_db"
      secret_name   = "delivery-mysql-cluster-secret"
      app_user      = "delivery_app"
      instances     = 1
      version       = "8.4.3"
      storage_size  = "10Gi"
    }
    notification = {
      cluster_name  = "notification-mysql"
      database_name = "notification_db"
      secret_name   = "notification-mysql-cluster-secret"
      app_user      = "notification_app"
      instances     = 1
      version       = "8.4.3"
      storage_size  = "10Gi"
    }
    support = {
      cluster_name  = "support-mysql"
      database_name = "support_db"
      secret_name   = "support-mysql-cluster-secret"
      app_user      = "support_app"
      instances     = 1
      version       = "8.4.3"
      storage_size  = "10Gi"
    }
    settlement = {
      cluster_name  = "settlement-mysql"
      database_name = "settlement_db"
      secret_name   = "settlement-mysql-cluster-secret"
      app_user      = "settlement_app"
      instances     = 1
      version       = "8.4.3"
      storage_size  = "10Gi"
    }
  }

  cross_service_readonly = [
    {
      name           = "settlement"
      username       = "settlement_ro"
      target_service = "payment"
      tables = ["payments", "refunds"]
    },
  ]
}

prevent_destroy = true

terraform {
  source = "${get_repo_root()}/infra//resource/mysql"
}
