terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.87.0"
    }
    harbor = {
      source  = "goharbor/harbor"
      version = "3.10.16"
    }
  }
}

provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

provider "vault" {
  address      = var.vault_url
  token        = var.vault_token
  ca_cert_file = var.kube_ca_crt_path
}

provider "aws" {
  region                      = "us-west-1"
  access_key                  = var.aws_access_key
  secret_key                  = var.aws_secret_key
  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    s3 = var.minio_api_url
  }
}

provider "harbor" {
  url      = var.harbor_url
  username = var.harbor_username
  password = var.harbor_password
}

resource "kubernetes_namespace" "harbor" {
  metadata {
    name = var.namespace
    labels = {
      "istio-injection" = "enabled"
    }
  }
}

module "storage" {
  source       = "./modules/storage"
  bucket_names = var.bucket_names
  providers = {
    aws = aws
  }
}

module "helm" {
  source         = "./modules/helm"
  namespace      = kubernetes_namespace.harbor.metadata.0.name
  admin_password = var.harbor_password
  aws_access_key = var.aws_access_key
  aws_secret_key = var.aws_secret_key
  providers = {
    helm = helm
  }
}

module "istio" {
  source       = "./modules/istio"
  namespace    = kubernetes_namespace.harbor.metadata.0.name
  host_address = var.harbor_address
  host_fqdn    = var.harbor_fqdn
  dest_port    = 80
  name_prefix  = "harbor"
}

module "user" {
  source       = "./modules/user"
  harbor_users = var.harbor_users
  providers = {
    harbor = harbor
  }
}

module "project" {
  source          = "./modules/project"
  harbor_projects = var.harbor_projects
}

module "cert" {
  source           = "./modules/cert"
  namespace        = kubernetes_namespace.harbor.metadata.0.name
  harbor_address   = var.harbor_address
  kube_ca_crt_path = var.kube_ca_crt_path
  providers = {
    vault = vault
  }
}
