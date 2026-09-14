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
  mount = "kv-infra"
  name  = "grafana/admin/credentials"
}

provider "grafana" {
  url  = var.grafana_url
  auth = "${ephemeral.vault_kv_secret_v2.grafana_admin.data["admin_user"]}:${ephemeral.vault_kv_secret_v2.grafana_admin.data["admin_password"]}"
}

resource "kubernetes_namespace_v1" "kyverno" {
  metadata {
    name = var.namespace
  }
}

module "helm" {
  source              = "./modules/helm"
  namespace           = kubernetes_namespace_v1.kyverno.metadata[0].name
  chart_version       = var.chart_version
  image_pull_secret   = var.image_pull_secret
  ca_bundle_configmap = var.ca_bundle_configmap
  ca_bundle_key       = var.ca_bundle_key
  providers = {
    helm = helm
  }
}

module "alertrules" {
  source            = "../../shared/grafana-alertrules"
  group_name        = "kyverno"
  folder_title      = "Kyverno"
  parent_folder_uid = "alerts"
  rules = [
    {
      name        = "KyvernoAdmissionDown"
      expr        = "absent(up{job=\"kyverno-svc-metrics\"} == 1)"
      for         = "5m"
      severity    = "critical"
      summary     = "Kyverno admission controller is down"
      description = "No healthy kyverno admission controller has been scraped for 5 minutes. The resource webhook fails open, so unsigned images are admitted while it is down."
    },
    {
      name        = "KyvernoAdmissionSlow"
      expr        = "histogram_quantile(0.99, sum by (le) (rate(kyverno_admission_review_duration_seconds_bucket[5m]))) > 5"
      for         = "10m"
      severity    = "warning"
      summary     = "Kyverno admission reviews are slow"
      description = "p99 admission review latency is {{ $value | humanizeDuration }}. The webhook fails open on timeout, so slow reviews skip image verification."
    },
    {
      name        = "KyvernoEnforceBlocked"
      expr        = "sum by (policy_name, resource_kind) (increase(kyverno_image_validating_policy_execution_duration_seconds_count{policy_validation_mode=\"enforce\", result=\"fail\"}[10m])) > 0"
      for         = "0m"
      severity    = "warning"
      summary     = "Kyverno rejected an image ({{ $labels.policy_name }})"
      description = "Policy {{ $labels.policy_name }} failed verification for {{ $value }} {{ $labels.resource_kind }} admission(s) in the last 10 minutes. Check the image signature of the deployment."
    },
    {
      name        = "KyvernoPolicyEvaluationErrors"
      expr        = "sum by (policy_name) (increase(kyverno_image_validating_policy_execution_duration_seconds_count{result=\"error\"}[10m])) > 0"
      for         = "0m"
      severity    = "warning"
      summary     = "Kyverno policy evaluation errors ({{ $labels.policy_name }})"
      description = "Policy {{ $labels.policy_name }} hit {{ $value }} evaluation error(s) in the last 10 minutes. With failurePolicy Ignore an error admits the resource without verification."
    },
  ]
  providers = {
    grafana = grafana
  }
}
