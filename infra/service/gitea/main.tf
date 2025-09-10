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
    grafana = {
      source  = "grafana/grafana"
      version = "~> 4.44"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.9"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.3"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.3"
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

resource "kubernetes_namespace_v1" "gitea" {
  metadata {
    name = var.namespace
    labels = {
      "istio.io/dataplane-mode" = "ambient"
    }
  }
}

module "secret" {
  source              = "./modules/secret"
  namespace           = kubernetes_namespace_v1.gitea.metadata[0].name
  postgresql_username = var.admin_username
  admin_username      = var.admin_username
  admin_email         = var.admin_email
  providers = {
    vault      = vault
    kubernetes = kubernetes
  }
}

module "helm" {
  source                   = "./modules/helm"
  namespace                = kubernetes_namespace_v1.gitea.metadata[0].name
  ssh_host_key_secret_name = module.secret.ssh_host_key_secret_name
  admin_email              = var.admin_email
  valkey_password          = module.secret.valkey_password
  postgresql_username      = module.secret.postgresql_username
  postgresql_password      = module.secret.postgresql_password
  providers = {
    helm = helm
  }
}

module "actions" {
  source             = "./modules/actions"
  namespace          = kubernetes_namespace_v1.gitea.metadata[0].name
  gitea_service_name = var.gitea_service_name
  gitea_api_url      = "https://${var.gitea_address}/api/v1"
  admin_username     = var.admin_username
  admin_password     = module.secret.admin_password
  providers = {
    helm = helm
  }
  depends_on = [module.helm]
}

module "route" {
  source          = "../../shared/gateway-api"
  hostname        = var.gitea_address
  service_name    = var.gitea_service_name
  service_port    = 3000
  route_namespace = kubernetes_namespace_v1.gitea.metadata[0].name
  name_prefix     = "gitea"
  request_timeout = "10m"
}

module "alertrules" {
  source            = "../../shared/grafana-alertrules"
  group_name        = "gitea"
  folder_title      = "Gitea"
  parent_folder_uid = "alerts"
  rules = [
    {
      name        = "GiteaDown"
      expr        = "up{job=~\".*gitea.*\"} == 0"
      for         = "5m"
      severity    = "critical"
      summary     = "Gitea instance is down"
      description = "Gitea instance {{ $labels.instance }} has been down for more than 5 minutes."
    },
    {
      name        = "GiteaHighResponseTime"
      expr        = "histogram_quantile(0.99, rate(gitea_http_request_duration_seconds_bucket[5m])) > 5"
      for         = "10m"
      severity    = "warning"
      summary     = "Gitea high response time"
      description = "Gitea instance {{ $labels.instance }} p99 response time is {{ $value }}s, exceeding 5s threshold."
    },
    {
      name        = "GiteaHighErrorRate"
      expr        = "rate(gitea_http_request_total{code=~\"5..\"}[5m]) / rate(gitea_http_request_total[5m]) > 0.05"
      for         = "10m"
      severity    = "warning"
      summary     = "Gitea high error rate"
      description = "Gitea instance {{ $labels.instance }} 5xx error rate is {{ $value | humanizePercentage }}, exceeding 5% threshold."
    },
  ]
  providers = {
    grafana = grafana
  }
}
