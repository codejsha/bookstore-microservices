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
    "../../service/postgresql",
    "../../service/vault",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role = "postgresql-config"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  namespace         = "bookstore"
  postgres_services = ["catalog", "customer", "identity", "inventory"]
  debezium_username = "debezium"

  postgres_db_map = {
    catalog   = "catalog_db"
    customer  = "customer_db"
    identity  = "identity_db"
    inventory = "inventory_db"
  }
  postgres_app_user_map = {
    catalog   = "catalog_app"
    customer  = "customer_app"
    identity  = "identity_app"
    inventory = "inventory_app"
  }

  postgres_cluster_config = {
    catalog = {
      cluster_name         = "catalog-postgres"
      image                = "ghcr.io/cloudnative-pg/postgresql:18"
      instances            = 1
      database             = "catalog_db"
      app_user             = "catalog"
      storage_size         = "20Gi"
      monitoring           = true
      cdc                  = true
      debezium_secret_name = "catalog-postgres-debezium"
      resources = {
        limits   = { cpu = "1", memory = "2Gi" }
        requests = { cpu = "100m", memory = "512Mi" }
      }
    }
    customer = {
      cluster_name = "customer-postgres"
      image        = "ghcr.io/cloudnative-pg/postgresql:18"
      instances    = 1
      database     = "customer_db"
      app_user     = "customer"
      storage_size = "10Gi"
      monitoring   = true
      resources = {
        limits   = { cpu = "1", memory = "1Gi" }
        requests = { cpu = "50m", memory = "256Mi" }
      }
    }
    identity = {
      cluster_name = "identity-postgres"
      image        = "ghcr.io/cloudnative-pg/postgresql:18"
      instances    = 1
      database     = "identity_db"
      app_user     = "identity"
      storage_size = "10Gi"
      monitoring   = true
      resources = {
        limits   = { cpu = "1", memory = "1Gi" }
        requests = { cpu = "50m", memory = "256Mi" }
      }
    }
    inventory = {
      cluster_name = "inventory-postgres"
      image        = "ghcr.io/cloudnative-pg/postgresql:18"
      instances    = 1
      database     = "inventory_db"
      app_user     = "inventory"
      storage_size = "10Gi"
      monitoring   = true
      resources = {
        limits   = { cpu = "1", memory = "1Gi" }
        requests = { cpu = "50m", memory = "256Mi" }
      }
    }
  }
}

prevent_destroy = true

terraform {
  source = "${get_repo_root()}/infra//resource/postgresql"
}
