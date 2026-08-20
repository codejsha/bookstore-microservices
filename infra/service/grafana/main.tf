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
    random = {
      source  = "hashicorp/random"
      version = "~> 3.9"
    }
  }
}

data "vault_kv_secret_v2" "seaweedfs" {
  mount = "kv"
  name  = "seaweedfs/s3/credentials"
}

resource "kubernetes_namespace_v1" "grafana" {
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
    namespace = kubernetes_namespace_v1.grafana.metadata[0].name
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

module "loki" {
  source           = "./modules/loki"
  namespace        = kubernetes_namespace_v1.grafana.metadata[0].name
  s3_endpoint      = var.s3_endpoint
  s3_region        = var.s3_region
  s3_bucket        = var.loki_s3_bucket
  s3_access_key_id = data.vault_kv_secret_v2.seaweedfs.data["access_key_id"]
  s3_secret_key    = data.vault_kv_secret_v2.seaweedfs.data["secret_access_key"]
  providers = {
    helm = helm
  }
}

module "tempo" {
  source           = "./modules/tempo"
  namespace        = kubernetes_namespace_v1.grafana.metadata[0].name
  s3_endpoint      = var.s3_endpoint
  s3_region        = var.s3_region
  s3_bucket        = var.tempo_s3_bucket
  s3_access_key_id = data.vault_kv_secret_v2.seaweedfs.data["access_key_id"]
  s3_secret_key    = data.vault_kv_secret_v2.seaweedfs.data["secret_access_key"]
  providers = {
    helm = helm
  }
}

module "pyroscope" {
  source           = "./modules/pyroscope"
  namespace        = kubernetes_namespace_v1.grafana.metadata[0].name
  s3_endpoint      = var.s3_endpoint
  s3_region        = var.s3_region
  s3_bucket        = var.pyroscope_s3_bucket
  s3_access_key_id = data.vault_kv_secret_v2.seaweedfs.data["access_key_id"]
  s3_secret_key    = data.vault_kv_secret_v2.seaweedfs.data["secret_access_key"]
  providers = {
    helm = helm
  }
}

module "alloy" {
  source              = "./modules/alloy"
  namespace           = kubernetes_namespace_v1.grafana.metadata[0].name
  gateway_target_fqdn = "alloy-gateway.${kubernetes_namespace_v1.grafana.metadata[0].name}.svc.cluster.local"
  loki_push_url       = "http://loki-gateway.${kubernetes_namespace_v1.grafana.metadata[0].name}.svc.cluster.local/otlp"
  tempo_otlp_endpoint = "tempo.${kubernetes_namespace_v1.grafana.metadata[0].name}.svc.cluster.local:4317"
  prometheus_otlp_url = "http://prometheus-kube-prometheus-prometheus.prometheus.svc.cluster.local:9090/api/v1/otlp"
  pyroscope_write_url = "http://pyroscope.${kubernetes_namespace_v1.grafana.metadata[0].name}.svc.cluster.local:4040"
  providers = {
    helm = helm
  }
  depends_on = [module.loki, module.tempo, module.pyroscope]
}

module "secret" {
  source    = "./modules/secret"
  namespace = kubernetes_namespace_v1.grafana.metadata[0].name
  providers = {
    vault      = vault
    kubernetes = kubernetes
  }
}

module "grafana" {
  source    = "./modules/grafana"
  namespace = kubernetes_namespace_v1.grafana.metadata[0].name
  providers = {
    helm = helm
  }
  depends_on = [module.loki, module.tempo, module.pyroscope, module.secret]
}

module "dashboard" {
  source    = "./modules/dashboard"
  namespace = kubernetes_namespace_v1.grafana.metadata[0].name
}

module "grafana_route" {
  source          = "../../shared/gateway-api"
  hostname        = var.grafana_address
  service_name    = var.grafana_service_name
  service_port    = 80
  route_namespace = kubernetes_namespace_v1.grafana.metadata[0].name
  name_prefix     = "grafana"
}

module "alloy_grpc_route" {
  source          = "../../shared/gateway-api"
  hostname        = var.alloy_grpc_address
  service_name    = var.alloy_gateway_service_name
  service_port    = 4317
  route_namespace = kubernetes_namespace_v1.grafana.metadata[0].name
  name_prefix     = "alloy-grpc"
  request_timeout = null
}

module "alloy_http_route" {
  source          = "../../shared/gateway-api"
  hostname        = var.alloy_http_address
  service_name    = var.alloy_gateway_service_name
  service_port    = 4318
  route_namespace = kubernetes_namespace_v1.grafana.metadata[0].name
  name_prefix     = "alloy-http"
}

resource "kubernetes_manifest" "alloy_gateway_reference_grant" {
  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1beta1"
    kind       = "ReferenceGrant"
    metadata = {
      name      = "alloy-gateway-from-app-routes"
      namespace = kubernetes_namespace_v1.grafana.metadata[0].name
    }
    spec = {
      from = [
        for ns in var.telemetry_route_namespaces : {
          group     = "gateway.networking.k8s.io"
          kind      = "HTTPRoute"
          namespace = ns
        }
      ]
      to = [
        {
          group = ""
          kind  = "Service"
          name  = var.alloy_gateway_service_name
        }
      ]
    }
  }
}
