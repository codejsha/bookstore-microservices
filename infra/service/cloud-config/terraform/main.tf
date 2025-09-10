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
    local = {
      source  = "hashicorp/local"
      version = "~> 2.9"
    }
    external = {
      source  = "hashicorp/external"
      version = "~> 2.4"
    }
  }
}

resource "kubernetes_namespace_v1" "config_server" {
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
    namespace = kubernetes_namespace_v1.config_server.metadata[0].name
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

data "kubernetes_secret_v1" "wildcard_ca" {
  metadata {
    name      = "wildcard-cert"
    namespace = "istio-system"
  }
}

resource "local_file" "wildcard_ca" {
  content  = data.kubernetes_secret_v1.wildcard_ca.data["ca.crt"]
  filename = "${path.module}/.wildcard-ca.crt"
}

data "external" "chart_rev" {
  program = ["bash", "-c",
    "printf '{\"sha\":\"%s\"}' \"$(git -C '${var.repo_root}' rev-parse --short=8 'HEAD:infra/service/cloud-config/helm' 2>/dev/null || echo none)\""
  ]
}

module "helm" {
  source             = "./modules/helm"
  namespace          = kubernetes_namespace_v1.config_server.metadata[0].name
  repository_ca_file = local_file.wildcard_ca.filename
  chart_version      = "1.0.0-${data.external.chart_rev.result.sha}"
  providers = {
    helm = helm
  }
}

resource "kubernetes_manifest" "config_server_authz" {
  count = length(var.client_principals) > 0 ? 1 : 0
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"
    metadata = {
      name      = "config-server-authz"
      namespace = kubernetes_namespace_v1.config_server.metadata[0].name
    }
    spec = {
      selector = {
        matchLabels = {
          "app.kubernetes.io/name" = "config-server"
        }
      }
      action = "ALLOW"
      rules = [
        {
          from = [
            {
              source = {
                principals = var.client_principals
              }
            }
          ]
        }
      ]
    }
  }
  depends_on = [module.helm]
}
