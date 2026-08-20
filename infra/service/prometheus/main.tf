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

resource "kubernetes_namespace_v1" "prometheus" {
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
    namespace = kubernetes_namespace_v1.prometheus.metadata[0].name
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
  namespace = kubernetes_namespace_v1.prometheus.metadata[0].name
  providers = {
    helm = helm
  }
}

module "prometheus_route" {
  source          = "../../shared/gateway-api"
  hostname        = var.prometheus_address
  service_name    = var.prometheus_service_name
  service_port    = 9090
  route_namespace = kubernetes_namespace_v1.prometheus.metadata[0].name
  name_prefix     = "prometheus"
  request_timeout = "10m"
}

module "servicemonitor" {
  source     = "./modules/servicemonitor"
  depends_on = [module.helm]
}

module "alertrules" {
  source            = "../../shared/grafana-alertrules"
  group_name        = "prometheus"
  folder_title      = "Prometheus"
  parent_folder_uid = "alerts"
  rules = [
    {
      name        = "IstioZtunnelDown"
      expr        = "up{job=~\".*ztunnel.*\"} == 0"
      for         = "5m"
      severity    = "critical"
      summary     = "Istio ztunnel down ({{ $labels.instance }})"
      description = "Istio ztunnel on {{ $labels.instance }} has been unreachable for 5 minutes. Ambient mesh traffic is disrupted."
    },
    {
      name        = "IstioPilotDown"
      expr        = "up{job=~\".*istiod.*\"} == 0"
      for         = "5m"
      severity    = "critical"
      summary     = "Istiod (pilot) down"
      description = "Istiod has been unreachable for 5 minutes. Service mesh configuration updates will fail."
    },
    {
      name        = "IstioPilotXdsPushErrors"
      expr        = "rate(pilot_xds_push_errors[5m]) > 0"
      for         = "5m"
      severity    = "warning"
      summary     = "Istio xDS push errors"
      description = "Istiod is experiencing xDS push errors at {{ $value }}/sec. Envoy configurations may be stale."
    },
    {
      name        = "IstioHighRequestErrorRate"
      expr        = "sum(rate(istio_requests_total{response_code=~\"5..\"}[5m])) / sum(rate(istio_requests_total[5m])) > 0.05"
      for         = "10m"
      severity    = "warning"
      summary     = "Istio mesh high 5xx error rate"
      description = "Overall mesh 5xx error rate exceeds 5% for 10 minutes."
    },
    {
      name        = "IstioHighRequestLatency"
      expr        = "histogram_quantile(0.99, sum(rate(istio_request_duration_milliseconds_bucket[5m])) by (le, destination_service)) > 5000"
      for         = "10m"
      severity    = "warning"
      summary     = "Istio high request latency ({{ $labels.destination_service }})"
      description = "Service {{ $labels.destination_service }} p99 latency exceeds 5 seconds for 10 minutes."
    },
    {
      name        = "IstioMTLSPolicyError"
      expr        = "sum(rate(istio_requests_total{connection_security_policy!=\"mutual_tls\",destination_service_namespace!=\"\"}[5m])) > 0"
      for         = "15m"
      severity    = "warning"
      summary     = "Non-mTLS traffic detected in mesh"
      description = "Traffic without mutual TLS detected for 15 minutes. Check PeerAuthentication policies."
    },
    {
      name        = "NodeCPUUsageHigh"
      expr        = "1 - avg by (instance) (rate(node_cpu_seconds_total{mode=\"idle\"}[5m])) > 0.9"
      for         = "15m"
      severity    = "warning"
      summary     = "Node CPU usage high ({{ $labels.instance }})"
      description = "Node {{ $labels.instance }} CPU usage is {{ $value | humanizePercentage }}, exceeding 90% for 15 minutes."
    },
    {
      name        = "NodeMemoryUsageHigh"
      expr        = "1 - node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes > 0.9"
      for         = "15m"
      severity    = "warning"
      summary     = "Node memory usage high ({{ $labels.instance }})"
      description = "Node {{ $labels.instance }} memory usage is {{ $value | humanizePercentage }}, exceeding 90% for 15 minutes."
    },
    {
      name        = "PodCPUUsageHigh"
      expr        = "sum by (namespace, pod, container) (rate(container_cpu_usage_seconds_total{container!=\"\",image!=\"\"}[5m])) / sum by (namespace, pod, container) (kube_pod_container_resource_limits{resource=\"cpu\"}) > 0.8"
      for         = "15m"
      severity    = "warning"
      summary     = "Pod CPU usage high ({{ $labels.namespace }}/{{ $labels.pod }})"
      description = "Container {{ $labels.container }} in pod {{ $labels.pod }} ({{ $labels.namespace }}) is using {{ $value | humanizePercentage }} of its CPU limit, exceeding 80% for 15 minutes."
    },
    {
      name        = "PodMemoryUsageHigh"
      expr        = "sum by (namespace, pod, container) (container_memory_working_set_bytes{container!=\"\",image!=\"\"}) / sum by (namespace, pod, container) (kube_pod_container_resource_limits{resource=\"memory\"}) > 0.8"
      for         = "15m"
      severity    = "warning"
      summary     = "Pod memory usage high ({{ $labels.namespace }}/{{ $labels.pod }})"
      description = "Container {{ $labels.container }} in pod {{ $labels.pod }} ({{ $labels.namespace }}) is using {{ $value | humanizePercentage }} of its memory limit, exceeding 80% for 15 minutes."
    },
    {
      name        = "PVCUsageHigh"
      expr        = "kubelet_volume_stats_used_bytes / kubelet_volume_stats_capacity_bytes > 0.85"
      for         = "15m"
      severity    = "warning"
      summary     = "PVC usage high ({{ $labels.persistentvolumeclaim }})"
      description = "PVC {{ $labels.persistentvolumeclaim }} in namespace {{ $labels.namespace }} is {{ $value | humanizePercentage }} full."
    },
    {
      name        = "PVCUsageCritical"
      expr        = "kubelet_volume_stats_used_bytes / kubelet_volume_stats_capacity_bytes > 0.95"
      for         = "5m"
      severity    = "critical"
      summary     = "PVC usage critical ({{ $labels.persistentvolumeclaim }})"
      description = "PVC {{ $labels.persistentvolumeclaim }} in namespace {{ $labels.namespace }} is {{ $value | humanizePercentage }} full. Immediate action required."
    },
    {
      name        = "PodCrashLooping"
      expr        = "increase(kube_pod_container_status_restarts_total[1h]) > 5"
      for         = "10m"
      severity    = "warning"
      summary     = "Pod crash looping ({{ $labels.namespace }}/{{ $labels.pod }})"
      description = "Pod {{ $labels.pod }} in namespace {{ $labels.namespace }} has restarted {{ $value }} times in the last hour."
    },
    {
      name        = "DeploymentReplicasMismatch"
      expr        = "kube_deployment_spec_replicas != kube_deployment_status_ready_replicas"
      for         = "15m"
      severity    = "warning"
      summary     = "Deployment replicas mismatch ({{ $labels.namespace }}/{{ $labels.deployment }})"
      description = "Deployment {{ $labels.deployment }} in {{ $labels.namespace }} has {{ $value }} ready vs desired replicas for 15 minutes."
    },
    {
      name        = "StatefulSetReplicasMismatch"
      expr        = "kube_statefulset_replicas != kube_statefulset_status_replicas_ready"
      for         = "15m"
      severity    = "warning"
      summary     = "StatefulSet replicas mismatch ({{ $labels.namespace }}/{{ $labels.statefulset }})"
      description = "StatefulSet {{ $labels.statefulset }} in {{ $labels.namespace }} has unavailable replicas for 15 minutes."
    },
    {
      name        = "ContainerOOMKilled"
      expr        = "kube_pod_container_status_last_terminated_reason{reason=\"OOMKilled\"} == 1"
      for         = "0m"
      severity    = "warning"
      summary     = "Container OOM killed ({{ $labels.namespace }}/{{ $labels.pod }})"
      description = "Container {{ $labels.container }} in pod {{ $labels.pod }} was OOM killed. Consider increasing memory limits."
    },
  ]
  providers = {
    grafana = grafana
  }
}

module "settlement_alerts" {
  source            = "../../shared/grafana-alertrules"
  group_name        = "settlement"
  folder_title      = "Settlement"
  parent_folder_uid = "alerts"
  rules = [
    {
      name        = "SettlementJobFailed"
      expr        = "kube_job_failed{namespace=\"bookstore\", condition=\"true\"} * on (namespace, job_name) group_left(owner_name) kube_job_owner{namespace=\"bookstore\", owner_name=\"settlement-batch\", owner_kind=\"CronJob\"} > 0"
      for         = "5m"
      severity    = "critical"
      summary     = "Settlement batch job failed ({{ $labels.job_name }})"
      description = "Settlement reconciliation Job {{ $labels.job_name }} (CronJob settlement-batch) reached the Failed condition. Daily settlement did not complete — a payment/PG discrepancy or extract failure is likely. Inspect the Job logs before the ttlSecondsAfterFinished window (24h) removes the pod."
    },
    {
      name        = "SettlementCronJobStale"
      expr        = "(time() - (max(kube_cronjob_status_last_successful_time{cronjob=\"settlement-batch\", namespace=\"bookstore\"}) or vector(0))) > 90000"
      for         = "15m"
      severity    = "critical"
      summary     = "Settlement CronJob has no recent successful run"
      description = "CronJob settlement-batch (namespace bookstore) has not completed successfully in over 25h (or has never succeeded). The daily 02:00 KST settlement is missing — check for a suspended CronJob, a startingDeadlineSeconds overrun, or repeated Job failures."
    },
  ]
  providers = {
    grafana = grafana
  }
}
