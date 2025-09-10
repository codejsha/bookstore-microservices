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

resource "kubernetes_namespace_v1" "flink" {
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
    namespace = kubernetes_namespace_v1.flink.metadata[0].name
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
  source    = "./modules/helm"
  namespace = kubernetes_namespace_v1.flink.metadata[0].name
  providers = {
    helm = helm
  }
}

module "session_cluster" {
  source     = "./modules/session-cluster"
  namespace  = kubernetes_namespace_v1.flink.metadata[0].name
  depends_on = [module.helm]
}

module "alertrules" {
  source            = "../../shared/grafana-alertrules"
  group_name        = "flink"
  folder_title      = "Flink"
  parent_folder_uid = "alerts"
  rules = [
    {
      name        = "FlinkJobFailed"
      expr        = "flink_jobmanager_job_uptime == 0 and flink_jobmanager_job_restartingTime > 0"
      for         = "5m"
      severity    = "critical"
      summary     = "Flink job has failed"
      description = "Flink job {{ $labels.job_name }} has failed and is not recovering on {{ $labels.instance }}."
    },
    {
      name        = "FlinkCheckpointFailing"
      expr        = "increase(flink_jobmanager_job_numberOfFailedCheckpoints[10m]) > 3"
      for         = "10m"
      severity    = "warning"
      summary     = "Flink checkpoints are failing"
      description = "Flink job {{ $labels.job_name }} has had {{ $value }} failed checkpoints in the last 10 minutes."
    },
    {
      name        = "FlinkTaskManagerLost"
      expr        = "flink_jobmanager_numRegisteredTaskManagers < 1"
      for         = "5m"
      severity    = "critical"
      summary     = "Flink TaskManager lost"
      description = "Flink cluster has {{ $value }} registered TaskManagers, expected at least 1."
    },
    {
      name        = "FlinkHighBackPressure"
      expr        = "flink_taskmanager_job_task_backPressuredTimeMsPerSecond > 500"
      for         = "10m"
      severity    = "warning"
      summary     = "Flink task experiencing high back pressure"
      description = "Flink task {{ $labels.task_name }} in job {{ $labels.job_name }} has {{ $value }}ms/s of back pressure time, exceeding 500ms/s threshold."
    },
  ]
  providers = {
    grafana = grafana
  }
}
