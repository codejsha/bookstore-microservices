terraform {
  required_providers {
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

resource "grafana_folder" "alerts" {
  title = "Alerts"
  uid   = "alerts"
}

module "vault_alerts" {
  source            = "../../shared/grafana-alertrules"
  group_name        = "vault"
  folder_title      = "Vault"
  parent_folder_uid = grafana_folder.alerts.uid
  rules = [
    {
      name        = "VaultDown"
      expr        = "up{job=~\".*vault.*\"} == 0"
      for         = "5m"
      severity    = "critical"
      summary     = "Vault instance is down"
      description = "Vault instance {{ $labels.instance }} has been down for more than 5 minutes."
    },
    {
      name        = "VaultSealed"
      expr        = "vault_core_unsealed == 0"
      for         = "1m"
      severity    = "critical"
      summary     = "Vault instance is sealed"
      description = "Vault instance {{ $labels.instance }} is sealed. Immediate action required to unseal."
    },
    {
      name        = "VaultHighResponseTime"
      expr        = "histogram_quantile(0.99, rate(vault_barrier_get_bucket[5m])) > 0.5"
      for         = "10m"
      severity    = "warning"
      summary     = "Vault high response time"
      description = "Vault instance {{ $labels.instance }} p99 barrier get latency is {{ $value }}s, exceeding 0.5s threshold."
    },
    {
      name        = "VaultHighErrorRate"
      expr        = "rate(vault_audit_log_request_failure[5m]) > 0"
      for         = "5m"
      severity    = "warning"
      summary     = "Vault audit log request failures detected"
      description = "Vault instance {{ $labels.instance }} is experiencing audit log request failures at a rate of {{ $value }}/s."
    },
    {
      name        = "VaultLeaderLost"
      expr        = "vault_core_active == 0"
      for         = "5m"
      severity    = "critical"
      summary     = "Vault has no active leader"
      description = "No active Vault leader detected on {{ $labels.instance }} for more than 5 minutes."
    },
    {
      name        = "VaultTokenExpiringSoon"
      expr        = "vault_token_lookup_ttl < 3600"
      for         = "5m"
      severity    = "warning"
      summary     = "Vault token expiring soon"
      description = "Vault token on {{ $labels.instance }} has a TTL of {{ $value }}s, which is less than 1 hour."
    },
    {
      name        = "VaultPKICertIssuanceFailure"
      expr        = "increase(vault_secret_kv_count{mount_point=\"pki_int\"}[10m]) == 0 and increase(vault_route_handle_request_count{mount_point=\"pki_int\"}[10m]) > 0"
      for         = "10m"
      severity    = "warning"
      summary     = "Vault PKI certificate issuance failure"
      description = "PKI certificate issuance requests are being made on {{ $labels.instance }} but no new certificates are being issued."
    },
  ]
  providers = {
    grafana = grafana
  }
}

module "seaweedfs_alerts" {
  source            = "../../shared/grafana-alertrules"
  group_name        = "seaweedfs"
  folder_title      = "SeaweedFS Alerts"
  parent_folder_uid = grafana_folder.alerts.uid
  rules = [
    {
      name        = "SeaweedFSMasterDown"
      expr        = "up{job=~\".*seaweedfs-master.*\"} == 0"
      for         = "5m"
      severity    = "critical"
      summary     = "SeaweedFS master is down"
      description = "SeaweedFS master {{ $labels.instance }} has been down for more than 5 minutes."
    },
    {
      name        = "SeaweedFSVolumeServerDown"
      expr        = "up{job=~\".*seaweedfs-volume.*\"} == 0"
      for         = "5m"
      severity    = "critical"
      summary     = "SeaweedFS volume server is down"
      description = "SeaweedFS volume server {{ $labels.instance }} has been down for more than 5 minutes."
    },
    {
      name        = "SeaweedFSFilerDown"
      expr        = "up{job=~\".*seaweedfs-filer.*\"} == 0"
      for         = "5m"
      severity    = "critical"
      summary     = "SeaweedFS filer is down"
      description = "SeaweedFS filer {{ $labels.instance }} has been down for more than 5 minutes."
    },
    {
      name        = "SeaweedFSDiskUsageHigh"
      expr        = "seaweedfs_disk_used_bytes / seaweedfs_disk_total_bytes > 0.85"
      for         = "10m"
      severity    = "warning"
      summary     = "SeaweedFS disk usage high"
      description = "SeaweedFS node {{ $labels.instance }} disk usage is {{ $value | humanizePercentage }}, exceeding 85% threshold."
    },
  ]
  providers = {
    grafana = grafana
  }
}

module "cert_manager_alerts" {
  source            = "../../shared/grafana-alertrules"
  group_name        = "cert-manager"
  folder_title      = "cert-manager"
  parent_folder_uid = grafana_folder.alerts.uid
  rules = [
    {
      name        = "CertManagerDown"
      expr        = "up{job=~\".*cert-manager.*\"} == 0"
      for         = "5m"
      severity    = "critical"
      summary     = "cert-manager is down"
      description = "cert-manager instance {{ $labels.instance }} has been down for more than 5 minutes."
    },
    {
      name        = "CertificateExpiringSoon"
      expr        = "(certmanager_certificate_expiration_timestamp_seconds - time()) < 604800"
      for         = "1h"
      severity    = "warning"
      summary     = "Certificate expiring soon"
      description = "Certificate {{ $labels.name }} in namespace {{ $labels.namespace }} will expire in less than 7 days."
    },
    {
      name        = "CertificateExpiryCritical"
      expr        = "(certmanager_certificate_expiration_timestamp_seconds - time()) < 86400"
      for         = "10m"
      severity    = "critical"
      summary     = "Certificate expiry critical"
      description = "Certificate {{ $labels.name }} in namespace {{ $labels.namespace }} will expire in less than 24 hours."
    },
    {
      name        = "CertificateNotReady"
      expr        = "certmanager_certificate_ready_status{condition=\"False\"} == 1"
      for         = "15m"
      severity    = "warning"
      summary     = "Certificate is not ready"
      description = "Certificate {{ $labels.name }} in namespace {{ $labels.namespace }} has been in a not-ready state for more than 15 minutes."
    },
  ]
  providers = {
    grafana = grafana
  }
}

locals {
  bookstore_access_logs = "{service_name=~\".+\"} | namespace=\"bookstore\" | json | __error__=\"\" | user_uid != \"\""
}

module "security_alerts" {
  source            = "../../shared/grafana-alertrules"
  group_name        = "user-anomaly"
  folder_title      = "Security"
  parent_folder_uid = grafana_folder.alerts.uid
  datasource_uid    = "loki"
  datasource_type   = "loki"
  rules = [
    {
      name        = "UserRequestRateHigh"
      expr        = "sum by (user_uid) (count_over_time(${local.bookstore_access_logs} [5m])) > 1500"
      for         = "5m"
      severity    = "warning"
      summary     = "User is issuing an abnormal request volume"
      description = "user_uid {{ $labels.user_uid }} made {{ $value }} backend requests in 5 minutes (> 5 req/s sustained); likely automation or scraping."
    },
    {
      name        = "UserClientErrorBurst"
      expr        = "sum by (user_uid) (count_over_time(${local.bookstore_access_logs} | status >= 400 | status < 500 [5m])) > 100"
      for         = "5m"
      severity    = "warning"
      summary     = "User is generating a burst of 4xx responses"
      description = "user_uid {{ $labels.user_uid }} received {{ $value }} 4xx responses in 5 minutes; probing, fuzzing or a broken client."
    },
    {
      name        = "UserAuthorizationDenied"
      expr        = "sum by (user_uid) (count_over_time(${local.bookstore_access_logs} | status = 403 [5m])) > 20"
      for         = "0s"
      severity    = "critical"
      summary     = "User repeatedly denied by authorization"
      description = "user_uid {{ $labels.user_uid }} hit {{ $value }} 403s in 5 minutes; likely enumerating other users' resources."
    },
    {
      name        = "UserOrderPlacementBurst"
      expr        = "sum by (user_uid) (count_over_time(${local.bookstore_access_logs} | app=\"order\" | method=\"POST\" | path=\"/api/v1/orders/saga/place\" [10m])) > 10"
      for         = "0s"
      severity    = "warning"
      summary     = "User placing orders at an abnormal rate"
      description = "user_uid {{ $labels.user_uid }} started {{ $value }} place-order sagas in 10 minutes."
    },
    {
      name        = "UserCancellationBurst"
      expr        = "sum by (user_uid) (count_over_time(${local.bookstore_access_logs} | app=\"order\" | method=\"POST\" | path=~\"/api/v1/orders/[^/]+/cancellation\" [1h])) > 5"
      for         = "0s"
      severity    = "warning"
      summary     = "User cancelling orders at an abnormal rate"
      description = "user_uid {{ $labels.user_uid }} requested {{ $value }} order cancellations in 1 hour; possible refund farming."
    },
    {
      name        = "UserReviewWriteBurst"
      expr        = "sum by (user_uid) (count_over_time(${local.bookstore_access_logs} | app=\"customer\" | method=\"POST\" | path=~\"/api/v1/customers/[^/]+/reviews\" [10m])) > 20"
      for         = "0s"
      severity    = "warning"
      summary     = "User posting reviews at an abnormal rate"
      description = "user_uid {{ $labels.user_uid }} posted {{ $value }} reviews in 10 minutes; likely review spam."
    },
    {
      name        = "ClientIpManyUsers"
      expr        = "count by (client_ip) (sum by (client_ip, user_uid) (count_over_time(${local.bookstore_access_logs} [15m]))) > 10"
      for         = "0s"
      severity    = "critical"
      summary     = "Single client IP acting as many users"
      description = "client_ip {{ $labels.client_ip }} authenticated as {{ $value }} distinct users in 15 minutes; credential stuffing or shared automation."
    },
    {
      name            = "EdgeRateLimitTripping"
      datasource_uid  = "prometheus"
      datasource_type = "prometheus"
      expr            = "sum by (destination_service) (rate(istio_requests_total{reporter=\"source\", response_code=\"429\"}[5m])) > 1"
      for             = "5m"
      severity        = "warning"
      summary         = "Edge rate limit is rejecting traffic"
      description     = "{{ $labels.destination_service }} is returning {{ $value }} 429/s for 5 minutes; either an attack is being throttled or the bucket is too small."
    },
    {
      name            = "AuthFailureSpike"
      datasource_uid  = "prometheus"
      datasource_type = "prometheus"
      expr            = "sum by (destination_service) (rate(istio_requests_total{reporter=\"destination\", destination_workload_namespace=\"bookstore\", response_code=~\"401|403\"}[5m])) > 5"
      for             = "5m"
      severity        = "warning"
      summary         = "Authentication/authorization failures spiking"
      description     = "{{ $labels.destination_service }} is rejecting {{ $value }} req/s with 401/403 for 5 minutes; credential stuffing or token replay."
    },
  ]
  providers = {
    grafana = grafana
  }
}
