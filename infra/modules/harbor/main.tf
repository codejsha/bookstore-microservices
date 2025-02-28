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
  source    = "./modules/helm"
  namespace = kubernetes_namespace.harbor.metadata.0.name
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
  name_prefix  = var.name_prefix
}

module "user" {
  source = "./modules/user"
  admin_password = var.harbor_password
  providers = {
    harbor = harbor
  }
}
