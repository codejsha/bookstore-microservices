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
  }
}

data "vault_pki_secret_backend_issuer" "pki_int" {
  backend    = "pki_int"
  issuer_ref = "default"
}

ephemeral "vault_kv_secret_v2" "grafana_admin" {
  mount = "kv"
  name  = "grafana/admin/credentials"
}

provider "grafana" {
  url  = var.grafana_url
  auth = "${ephemeral.vault_kv_secret_v2.grafana_admin.data["admin_user"]}:${ephemeral.vault_kv_secret_v2.grafana_admin.data["admin_password"]}"
}

resource "kubernetes_namespace_v1" "youtrack" {
  metadata {
    name = var.namespace
    labels = {
      "istio.io/dataplane-mode" = "ambient"
    }
  }
}

module "secret" {
  source      = "./modules/secret"
  namespace   = kubernetes_namespace_v1.youtrack.metadata[0].name
  admin_email = var.admin_email
  providers = {
    vault = vault
  }
}

module "helm" {
  source              = "./modules/helm"
  namespace           = kubernetes_namespace_v1.youtrack.metadata[0].name
  chart_version       = var.chart_version
  youtrack_version    = var.youtrack_version
  youtrack_address    = var.youtrack_address
  storage_class       = var.storage_class
  data_storage_size   = var.data_storage_size
  logs_storage_size   = var.logs_storage_size
  backup_storage_size = var.backup_storage_size
  providers = {
    helm = helm
  }
}

module "cert" {
  source           = "./modules/cert"
  namespace        = kubernetes_namespace_v1.youtrack.metadata[0].name
  youtrack_address = var.youtrack_address
  kube_ca_cert     = data.vault_pki_secret_backend_issuer.pki_int.certificate
  providers = {
    vault = vault
  }
}

module "gateway_api" {
  source          = "../../shared/gateway-api"
  hostname        = var.youtrack_address
  service_name    = "youtrack"
  service_port    = 8080
  route_namespace = kubernetes_namespace_v1.youtrack.metadata[0].name
  name_prefix     = "youtrack"
  request_timeout = "10m"
}

module "alertrules" {
  source            = "../../shared/grafana-alertrules"
  group_name        = "youtrack"
  folder_title      = "YouTrack"
  parent_folder_uid = "alerts"
  rules = [
    {
      name        = "YouTrackDown"
      expr        = "up{job=~\".*youtrack.*\"} == 0"
      for         = "5m"
      severity    = "critical"
      summary     = "YouTrack instance is down"
      description = "YouTrack instance {{ $labels.instance }} has been down for more than 5 minutes."
    },
    {
      name        = "YouTrackHighResponseTime"
      expr        = "histogram_quantile(0.99, rate(http_request_duration_seconds_bucket{job=~\".*youtrack.*\"}[5m])) > 5"
      for         = "10m"
      severity    = "warning"
      summary     = "YouTrack high response time"
      description = "YouTrack instance {{ $labels.instance }} p99 response time is {{ $value }}s, exceeding 5s threshold."
    },
    {
      name        = "YouTrackHighErrorRate"
      expr        = "rate(http_request_total{job=~\".*youtrack.*\",code=~\"5..\"}[5m]) / rate(http_request_total{job=~\".*youtrack.*\"}[5m]) > 0.05"
      for         = "10m"
      severity    = "warning"
      summary     = "YouTrack high error rate"
      description = "YouTrack instance {{ $labels.instance }} 5xx error rate is {{ $value | humanizePercentage }}, exceeding 5% threshold."
    },
    {
      name        = "YouTrackDataPvcAlmostFull"
      expr        = "kubelet_volume_stats_available_bytes{persistentvolumeclaim=~\"youtrack-data.*\"} / kubelet_volume_stats_capacity_bytes{persistentvolumeclaim=~\"youtrack-data.*\"} < 0.10"
      for         = "10m"
      severity    = "warning"
      summary     = "YouTrack data PVC < 10% free"
      description = "PVC {{ $labels.persistentvolumeclaim }} on YouTrack has less than 10% free space remaining."
    },
  ]
  providers = {
    grafana = grafana
  }
}
