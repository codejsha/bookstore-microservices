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
  }
}

resource "kubernetes_namespace_v1" "hyperswitch" {
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
    namespace = kubernetes_namespace_v1.hyperswitch.metadata[0].name
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

module "secret" {
  source    = "./modules/secret"
  namespace = kubernetes_namespace_v1.hyperswitch.metadata[0].name
  providers = {
    vault      = vault
    kubernetes = kubernetes
  }
}

module "helm" {
  source     = "./modules/helm"
  namespace  = kubernetes_namespace_v1.hyperswitch.metadata[0].name
  depends_on = [module.secret]
  providers = {
    helm = helm
  }
}

resource "kubernetes_manifest" "waypoint" {
  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "Gateway"
    metadata = {
      name      = "hyperswitch-waypoint"
      namespace = kubernetes_namespace_v1.hyperswitch.metadata[0].name
      labels = {
        "istio.io/waypoint-for" = "service"
      }
    }
    spec = {
      gatewayClassName = "istio-waypoint"
      listeners = [
        {
          name     = "mesh"
          port     = 15008
          protocol = "HBONE"
        }
      ]
    }
  }
}

resource "kubernetes_labels" "server_use_waypoint" {
  depends_on  = [module.helm]
  api_version = "v1"
  kind        = "Service"
  metadata {
    name      = "hyperswitch-hyperswitch-server"
    namespace = kubernetes_namespace_v1.hyperswitch.metadata[0].name
  }
  labels = {
    "istio.io/use-waypoint"         = "hyperswitch-waypoint"
    "istio.io/ingress-use-waypoint" = "true"
  }
}

resource "kubernetes_manifest" "server_timeout_route" {
  depends_on = [module.helm]
  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "HTTPRoute"
    metadata = {
      name      = "hyperswitch-server-timeout"
      namespace = kubernetes_namespace_v1.hyperswitch.metadata[0].name
    }
    spec = {
      parentRefs = [
        {
          kind  = "Service"
          group = ""
          name  = "hyperswitch-hyperswitch-server"
        }
      ]
      rules = [
        {
          backendRefs = [
            {
              name = "hyperswitch-hyperswitch-server"
              port = 80
            }
          ]
          timeouts = {
            request = var.request_timeout
          }
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "hyperswitch_dashboard_route" {
  depends_on = [module.helm]
  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "HTTPRoute"
    metadata = {
      name      = "hyperswitch-route"
      namespace = kubernetes_namespace_v1.hyperswitch.metadata[0].name
    }
    spec = {
      hostnames = ["hyperswitch.example.com"]
      parentRefs = [{
        group       = "gateway.networking.k8s.io"
        kind        = "Gateway"
        name        = "bookstore-gateway"
        namespace   = "istio-system"
        sectionName = "https-hyperswitch"
      }]
      rules = [{
        matches = [{
          path = {
            type  = "PathPrefix"
            value = "/"
          }
        }]
        backendRefs = [{
          group  = ""
          kind   = "Service"
          name   = "hyperswitch-control-center"
          port   = 80
          weight = 1
        }]
      }]
    }
  }
}

resource "kubernetes_manifest" "hyperswitch_api_route" {
  depends_on = [module.helm]
  manifest = {
    apiVersion = "gateway.networking.k8s.io/v1"
    kind       = "HTTPRoute"
    metadata = {
      name      = "hyperswitch-api-route"
      namespace = kubernetes_namespace_v1.hyperswitch.metadata[0].name
    }
    spec = {
      hostnames = ["hyperswitch-api.example.com"]
      parentRefs = [{
        group       = "gateway.networking.k8s.io"
        kind        = "Gateway"
        name        = "bookstore-gateway"
        namespace   = "istio-system"
        sectionName = "https-hyperswitch-api"
      }]
      rules = [{
        matches = [{
          path = {
            type  = "PathPrefix"
            value = "/"
          }
        }]
        backendRefs = [{
          group  = ""
          kind   = "Service"
          name   = "hyperswitch-hyperswitch-server"
          port   = 80
          weight = 1
        }]
      }]
    }
  }
}
