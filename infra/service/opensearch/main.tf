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

resource "kubernetes_namespace_v1" "opensearch" {
  metadata {
    name = var.namespace
    labels = {
      "istio.io/dataplane-mode" = "ambient"
    }
  }
}

module "secret" {
  source    = "./modules/secret"
  namespace = var.namespace
  providers = {
    vault      = vault
    kubernetes = kubernetes
  }
}

module "helm" {
  source    = "./modules/helm"
  namespace = kubernetes_namespace_v1.opensearch.metadata[0].name
  providers = {
    helm = helm
  }
}

module "route" {
  source          = "../../shared/gateway-api"
  hostname        = var.opensearch_dashboards_address
  service_name    = var.opensearch_dashboards_service_name
  service_port    = 5601
  route_namespace = kubernetes_namespace_v1.opensearch.metadata[0].name
  name_prefix     = "opensearch-dashboards"
  request_timeout = "10m"
}

module "route_api" {
  source                 = "../../shared/gateway-api"
  hostname               = "opensearch-api.example.com"
  service_name           = "opensearch-cluster-master"
  service_port           = 9200
  route_namespace        = kubernetes_namespace_v1.opensearch.metadata[0].name
  name_prefix            = "opensearch-api"
  remove_request_headers = ["x-request-id"]
  request_timeout        = "10m"
}

module "alertrules" {
  source            = "../../shared/grafana-alertrules"
  group_name        = "opensearch"
  folder_title      = "OpenSearch"
  parent_folder_uid = "alerts"
  rules = [
    {
      name        = "OpenSearchDown"
      expr        = "up{job=~\".*opensearch.*\"} == 0"
      for         = "5m"
      severity    = "critical"
      summary     = "OpenSearch instance is down"
      description = "OpenSearch instance {{ $labels.instance }} has been down for more than 5 minutes."
    },
    {
      name        = "OpenSearchClusterRed"
      expr        = "opensearch_cluster_status == 2"
      for         = "5m"
      severity    = "critical"
      summary     = "OpenSearch cluster status is red"
      description = "OpenSearch cluster {{ $labels.cluster }} status has been red for more than 5 minutes."
    },
    {
      name        = "OpenSearchClusterYellow"
      expr        = "opensearch_cluster_status == 1"
      for         = "30m"
      severity    = "warning"
      summary     = "OpenSearch cluster status is yellow"
      description = "OpenSearch cluster {{ $labels.cluster }} status has been yellow for more than 30 minutes."
    },
    {
      name        = "OpenSearchDiskWatermarkHigh"
      expr        = "opensearch_fs_path_total_bytes > 0 and (opensearch_fs_path_total_bytes - opensearch_fs_path_available_bytes) / opensearch_fs_path_total_bytes > 0.85"
      for         = "10m"
      severity    = "warning"
      summary     = "OpenSearch disk usage high"
      description = "OpenSearch node {{ $labels.node }} disk usage is {{ $value | humanizePercentage }}, exceeding 85% watermark."
    },
    {
      name        = "OpenSearchDiskWatermarkCritical"
      expr        = "opensearch_fs_path_total_bytes > 0 and (opensearch_fs_path_total_bytes - opensearch_fs_path_available_bytes) / opensearch_fs_path_total_bytes > 0.95"
      for         = "5m"
      severity    = "critical"
      summary     = "OpenSearch disk usage critical"
      description = "OpenSearch node {{ $labels.node }} disk usage is {{ $value | humanizePercentage }}, exceeding 95% critical watermark."
    },
    {
      name        = "OpenSearchHeapUsageHigh"
      expr        = "opensearch_jvm_mem_heap_used_percent / 100 > 0.9"
      for         = "10m"
      severity    = "warning"
      summary     = "OpenSearch JVM heap usage high"
      description = "OpenSearch node {{ $labels.node }} JVM heap usage is {{ $value | humanizePercentage }}, exceeding 90% threshold."
    },
    {
      name        = "OpenSearchPendingTasksHigh"
      expr        = "opensearch_cluster_pending_tasks_number > 50"
      for         = "10m"
      severity    = "warning"
      summary     = "OpenSearch pending tasks high"
      description = "OpenSearch cluster {{ $labels.cluster }} has {{ $value }} pending tasks, exceeding threshold of 50."
    },
    {
      name        = "OpenSearchIndexingRejections"
      expr        = "rate(opensearch_threadpool_threads_rejected{type=\"write\"}[5m]) > 0"
      for         = "5m"
      severity    = "warning"
      summary     = "OpenSearch indexing rejections detected"
      description = "OpenSearch node {{ $labels.node }} is rejecting write requests at a rate of {{ $value }}/s."
    },
    {
      name        = "OpenSearchGCPauseHigh"
      expr        = "avg_over_time(opensearch_jvm_gc_collection_time_seconds[10m]) > 0.5"
      for         = "10m"
      severity    = "warning"
      summary     = "OpenSearch GC pause time high"
      description = "OpenSearch node {{ $labels.node }} average GC pause time is {{ $value }}s, exceeding 0.5s threshold."
    },
  ]
  providers = {
    grafana = grafana
  }
}
