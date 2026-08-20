terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 3.2"
    }
    vault = {
      source  = "hashicorp/vault"
      version = "~> 5.10"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.58"
    }
  }
}

provider "vault" {
  address = var.vault_url
  auth_login {
    path = "auth/approle/login"
    parameters = {
      role_id   = var.vault_role_id
      secret_id = var.vault_secret_id
    }
  }
}

provider "aws" {
  region                      = "us-east-1"
  access_key                  = var.s3_access_key
  secret_key                  = var.s3_secret_key
  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
  endpoints {
    s3 = var.s3_endpoint
  }
}

resource "kubernetes_namespace_v1" "seaweedfs" {
  metadata {
    name = var.namespace
    labels = {
      "istio.io/dataplane-mode" = "ambient"
    }
  }
}

resource "kubernetes_limit_range_v1" "resource_limits" {
  metadata {
    name      = "resource-limits"
    namespace = kubernetes_namespace_v1.seaweedfs.metadata[0].name
  }
  spec {
    limit {
      type = "Container"
      default_request = {
        cpu    = "10m"
        memory = "32Mi"
      }
    }
  }
}

module "helm" {
  source               = "./modules/helm"
  namespace            = kubernetes_namespace_v1.seaweedfs.metadata[0].name
  seaweedfs_access_key = var.s3_access_key
  seaweedfs_secret_key = var.s3_secret_key
  providers = {
    helm = helm
  }
}

module "storage" {
  source       = "./modules/storage"
  bucket_names = var.bucket_names
  providers = {
    aws = aws
  }
  depends_on = [module.helm, module.api_route]
}

module "admin_route" {
  source          = "../../shared/gateway-api"
  hostname        = var.seaweedfs_address
  service_name    = var.seaweedfs_service_name
  service_port    = 23646
  route_namespace = kubernetes_namespace_v1.seaweedfs.metadata[0].name
  name_prefix     = "seaweedfs-admin"
  request_timeout = "10m"
}

module "api_route" {
  source          = "../../shared/gateway-api"
  hostname        = var.seaweedfs_api_address
  service_name    = var.seaweedfs_api_service_name
  service_port    = 8333
  route_namespace = kubernetes_namespace_v1.seaweedfs.metadata[0].name
  name_prefix     = "seaweedfs-api"
  request_timeout = "10m"
}
