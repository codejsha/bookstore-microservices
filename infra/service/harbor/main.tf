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
    grafana = {
      source  = "grafana/grafana"
      version = "~> 4.44"
    }
  }
}

ephemeral "vault_kv_secret_v2" "grafana_admin" {
  mount = "kv"
  name  = "grafana/admin/credentials"
}

provider "grafana" {
  url  = var.grafana_url
  auth = "${ephemeral.vault_kv_secret_v2.grafana_admin.data["admin_user"]}:${ephemeral.vault_kv_secret_v2.grafana_admin.data["admin_password"]}"
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
    s3 = var.aws_s3_api_url
  }
}

resource "kubernetes_namespace_v1" "harbor" {
  metadata {
    name = var.namespace
    labels = {
      "istio.io/dataplane-mode" = "ambient"
    }
  }
}

module "secret" {
  source    = "./modules/secret"
  namespace = kubernetes_namespace_v1.harbor.metadata[0].name
  providers = {
    vault      = vault
    kubernetes = kubernetes
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
  namespace      = kubernetes_namespace_v1.harbor.metadata[0].name
  aws_access_key = var.aws_access_key
  aws_secret_key = var.aws_secret_key
  s3_bucket      = var.bucket_names[0]
  s3_endpoint    = var.s3_registry_endpoint
  providers = {
    helm = helm
  }
}

module "route" {
  source          = "../../shared/gateway-api"
  hostname        = var.harbor_address
  service_name    = var.harbor_service_name
  service_port    = 80
  route_namespace = kubernetes_namespace_v1.harbor.metadata[0].name
  name_prefix     = "harbor"
  request_timeout = "10m"
}

module "alertrules" {
  source            = "../../shared/grafana-alertrules"
  group_name        = "harbor"
  folder_title      = "Harbor"
  parent_folder_uid = "alerts"
  rules = [
    {
      name        = "HarborDown"
      expr        = "up{job=~\".*harbor.*\"} == 0"
      for         = "5m"
      severity    = "critical"
      summary     = "Harbor instance is down"
      description = "Harbor instance {{ $labels.instance }} has been down for more than 5 minutes."
    },
    {
      name        = "HarborHighErrorRate"
      expr        = "rate(harbor_core_http_request_total{code=~\"5..\"}[5m]) / rate(harbor_core_http_request_total[5m]) > 0.05"
      for         = "10m"
      severity    = "warning"
      summary     = "Harbor high error rate"
      description = "Harbor instance {{ $labels.instance }} 5xx error rate is {{ $value | humanizePercentage }}, exceeding 5% threshold."
    },
    {
      name        = "HarborStorageUsageHigh"
      expr        = "harbor_project_quota_usage_byte / harbor_project_quota_byte > 0.9"
      for         = "30m"
      severity    = "warning"
      summary     = "Harbor storage usage high"
      description = "Harbor project {{ $labels.project_name }} storage quota usage is {{ $value | humanizePercentage }}, exceeding 90% threshold."
    },
    {
      name        = "HarborComponentUnhealthy"
      expr        = "harbor_health == 0"
      for         = "5m"
      severity    = "critical"
      summary     = "Harbor component is unhealthy"
      description = "Harbor component {{ $labels.component }} on {{ $labels.instance }} has been unhealthy for more than 5 minutes."
    },
  ]
  providers = {
    grafana = grafana
  }
}
