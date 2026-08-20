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

resource "kubernetes_namespace_v1" "keycloak" {
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
    namespace = kubernetes_namespace_v1.keycloak.metadata[0].name
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

module "postgresql_component" {
  source    = "./modules/postgresql-component"
  namespace = kubernetes_namespace_v1.keycloak.metadata[0].name
}

module "keycloak_operator" {
  source = "./modules/keycloak-operator"
}

module "keycloak_component" {
  source           = "./modules/keycloak-component"
  namespace        = kubernetes_namespace_v1.keycloak.metadata[0].name
  keycloak_address = var.keycloak_address
  depends_on       = [module.keycloak_operator, module.postgresql_component]
}

module "secret" {
  source     = "./modules/secret"
  namespace  = kubernetes_namespace_v1.keycloak.metadata[0].name
  depends_on = [module.keycloak_component]
  providers = {
    kubernetes = kubernetes
    vault      = vault
  }
}

module "route" {
  source          = "../../shared/gateway-api"
  hostname        = var.keycloak_address
  service_name    = "keycloak-service"
  service_port    = 8080
  route_namespace = kubernetes_namespace_v1.keycloak.metadata[0].name
  name_prefix     = "keycloak"
  request_timeout = "10m"
}

module "alertrules" {
  source            = "../../shared/grafana-alertrules"
  group_name        = "keycloak"
  folder_title      = "Keycloak"
  parent_folder_uid = "alerts"
  rules = [
    {
      name        = "KeycloakDown"
      expr        = "up{job=~\".*keycloak.*\"} == 0"
      for         = "5m"
      severity    = "critical"
      summary     = "Keycloak instance is down"
      description = "Keycloak instance {{ $labels.instance }} has been down for more than 5 minutes."
    },
    {
      name        = "KeycloakHighLoginFailureRate"
      expr        = "rate(keycloak_login_error_total[5m]) / (rate(keycloak_login_total[5m]) + rate(keycloak_login_error_total[5m])) > 0.1"
      for         = "10m"
      severity    = "warning"
      summary     = "Keycloak high login failure rate"
      description = "Keycloak realm {{ $labels.realm }} login failure rate is {{ $value | humanizePercentage }}, exceeding 10% threshold."
    },
    {
      name        = "KeycloakHighResponseTime"
      expr        = "histogram_quantile(0.99, rate(keycloak_request_duration_seconds_bucket[5m])) > 2"
      for         = "10m"
      severity    = "warning"
      summary     = "Keycloak high response time"
      description = "Keycloak instance {{ $labels.instance }} p99 response time is {{ $value }}s, exceeding 2s threshold."
    },
  ]
  providers = {
    grafana = grafana
  }
}
