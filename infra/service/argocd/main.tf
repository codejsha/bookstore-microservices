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

ephemeral "vault_kv_secret_v2" "grafana_admin" {
  mount = "kv"
  name  = "grafana/admin/credentials"
}

data "vault_kv_secret_v2" "argocd_admin" {
  mount = "kv"
  name  = "argocd/admin/credentials"
}

provider "grafana" {
  url  = var.grafana_url
  auth = "${ephemeral.vault_kv_secret_v2.grafana_admin.data["admin_user"]}:${ephemeral.vault_kv_secret_v2.grafana_admin.data["admin_password"]}"
}

resource "kubernetes_namespace_v1" "argocd" {
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
    namespace = kubernetes_namespace_v1.argocd.metadata[0].name
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

data "vault_kv_secret_v2" "gitea_ssh_host" {
  mount = "kv"
  name  = "gitea/ssh/host"
}

module "helm" {
  source                = "./modules/helm"
  namespace             = kubernetes_namespace_v1.argocd.metadata[0].name
  admin_password_bcrypt = data.vault_kv_secret_v2.argocd_admin.data["password_bcrypt"]
  admin_password_mtime  = data.vault_kv_secret_v2.argocd_admin.data["password_mtime"]
  ssh_extra_hosts       = "${var.gitea_ssh_fqdn} ${data.vault_kv_secret_v2.gitea_ssh_host.data["public"]}"
  providers = {
    helm = helm
  }
}

module "route" {
  source          = "../../shared/gateway-api"
  hostname        = var.argocd_address
  service_name    = var.argocd_service_name
  service_port    = 80
  route_namespace = kubernetes_namespace_v1.argocd.metadata[0].name
  name_prefix     = "argocd"
  request_timeout = "10m"
}

module "alertrules" {
  source            = "../../shared/grafana-alertrules"
  group_name        = "argocd"
  folder_title      = "Argo CD"
  parent_folder_uid = "alerts"
  rules = [
    {
      name        = "ArgocdDown"
      expr        = "up{job=~\".*argocd.*\"} == 0"
      for         = "5m"
      severity    = "critical"
      summary     = "ArgoCD instance is down"
      description = "ArgoCD instance {{ $labels.instance }} has been down for more than 5 minutes."
    },
    {
      name        = "ArgocdAppSyncFailed"
      expr        = "argocd_app_info{sync_status=\"OutOfSync\"} == 1"
      for         = "15m"
      severity    = "warning"
      summary     = "ArgoCD application sync failed"
      description = "ArgoCD application {{ $labels.name }} in project {{ $labels.project }} has been OutOfSync for more than 15 minutes."
    },
    {
      name        = "ArgocdAppHealthDegraded"
      expr        = "argocd_app_info{health_status=~\"Degraded|Missing|Unknown\"} == 1"
      for         = "10m"
      severity    = "warning"
      summary     = "ArgoCD application health degraded"
      description = "ArgoCD application {{ $labels.name }} has health status {{ $labels.health_status }} for more than 10 minutes."
    },
    {
      name        = "ArgocdSyncError"
      expr        = "increase(argocd_app_sync_total{phase!=\"Succeeded\"}[5m]) > 0"
      for         = "5m"
      severity    = "warning"
      summary     = "ArgoCD sync errors detected"
      description = "ArgoCD application {{ $labels.name }} has experienced {{ $value }} failed sync attempts in the last 5 minutes."
    },
    {
      name        = "ArgocdHighReconciliationTime"
      expr        = "histogram_quantile(0.99, rate(argocd_app_reconcile_bucket[10m])) > 60"
      for         = "10m"
      severity    = "warning"
      summary     = "ArgoCD high reconciliation time"
      description = "ArgoCD p99 reconciliation time is {{ $value }}s, exceeding 60s threshold."
    },
  ]
  providers = {
    grafana = grafana
  }
}
