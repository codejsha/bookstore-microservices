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

resource "kubernetes_namespace_v1" "nexus" {
  metadata {
    name = var.namespace
    labels = {
      "istio.io/dataplane-mode" = "ambient"
    }
  }
}

module "secret" {
  source    = "./modules/secret"
  namespace = kubernetes_namespace_v1.nexus.metadata[0].name
  providers = {
    vault      = vault
    kubernetes = kubernetes
  }
}

module "helm" {
  source    = "./modules/helm"
  namespace = kubernetes_namespace_v1.nexus.metadata[0].name
  providers = {
    helm       = helm
    kubernetes = kubernetes
  }
}

module "route" {
  source          = "../../shared/gateway-api"
  hostname        = var.nexus_address
  service_name    = var.nexus_service_name
  service_port    = 8081
  route_namespace = kubernetes_namespace_v1.nexus.metadata[0].name
  name_prefix     = "nexus"
  request_timeout = "10m"
}

module "alertrules" {
  source            = "../../shared/grafana-alertrules"
  group_name        = "nexus"
  folder_title      = "Nexus"
  parent_folder_uid = "alerts"
  rules = [
    {
      name        = "NexusDown"
      expr        = "up{job=~\".*nexus.*\"} == 0"
      for         = "5m"
      severity    = "critical"
      summary     = "Nexus instance is down"
      description = "Nexus instance {{ $labels.instance }} has been down for more than 5 minutes."
    },
  ]
  providers = {
    grafana = grafana
  }
}
