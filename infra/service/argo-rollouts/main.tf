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

provider "grafana" {
  url  = var.grafana_url
  auth = "${ephemeral.vault_kv_secret_v2.grafana_admin.data["admin_user"]}:${ephemeral.vault_kv_secret_v2.grafana_admin.data["admin_password"]}"
}

resource "kubernetes_namespace_v1" "argo_rollouts" {
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
    namespace = kubernetes_namespace_v1.argo_rollouts.metadata[0].name
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
  source                = "./modules/helm"
  namespace             = kubernetes_namespace_v1.argo_rollouts.metadata[0].name
  gatewayapi_plugin_url = var.gatewayapi_plugin_url
  providers = {
    helm = helm
  }
}

module "alertrules" {
  source            = "../../shared/grafana-alertrules"
  group_name        = "argo-rollouts"
  folder_title      = "Argo Rollouts"
  parent_folder_uid = "alerts"
  rules = [
    {
      name        = "ArgoRolloutsDown"
      expr        = "up{job=~\".*argo-rollouts.*\"} == 0"
      for         = "5m"
      severity    = "critical"
      summary     = "Argo Rollouts controller is down"
      description = "Argo Rollouts controller {{ $labels.instance }} has been down for more than 5 minutes; canary steps and analysis stop progressing."
    },
    {
      name        = "ArgoRolloutsDegraded"
      expr        = "rollout_info{phase=\"Degraded\"} == 1"
      for         = "10m"
      severity    = "critical"
      summary     = "Rollout is degraded"
      description = "Rollout {{ $labels.name }} in {{ $labels.namespace }} has been Degraded for more than 10 minutes."
    },
    {
      name        = "ArgoRolloutsStuckPaused"
      expr        = "rollout_info{phase=\"Paused\"} == 1"
      for         = "1h"
      severity    = "warning"
      summary     = "Rollout paused for over an hour"
      description = "Rollout {{ $labels.name }} in {{ $labels.namespace }} has been Paused for more than 1 hour; a canary is waiting on promotion."
    },
    {
      name        = "ArgoRolloutsAnalysisRunFailing"
      expr        = "increase(analysis_run_metric_phase{phase=\"Failed\"}[10m]) > 0"
      for         = "5m"
      severity    = "warning"
      summary     = "Rollout analysis metric failing"
      description = "AnalysisRun metric {{ $labels.metric }} for {{ $labels.name }} in {{ $labels.namespace }} reported Failed in the last 10 minutes."
    },
  ]
  providers = {
    grafana = grafana
  }
}
