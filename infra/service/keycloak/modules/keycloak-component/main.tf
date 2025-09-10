resource "kubernetes_manifest" "keycloak" {
  manifest = {
    apiVersion = "k8s.keycloak.org/v2alpha1"
    kind       = "Keycloak"
    metadata = {
      name      = "keycloak"
      namespace = var.namespace
    }
    spec = {
      instances = 2
      resources = {
        requests = { cpu = "50m", memory = "768Mi" }
        limits   = { cpu = "1", memory = "1536Mi" }
      }
      db = {
        vendor   = "postgres"
        host     = "keycloak-db-rw"
        database = "postgres"
        usernameSecret = {
          name = "keycloak-db-superuser"
          key  = "username"
        }
        passwordSecret = {
          name = "keycloak-db-superuser"
          key  = "password"
        }
      },
      additionalOptions = [
        {
          name  = "metrics-enabled"
          value = "true"
        },
        {
          name  = "cache-embedded-mtls-enabled"
          value = "false"
        },
        {
          name  = "proxy-headers"
          value = "xforwarded"
        }
      ]
      http = {
        httpEnabled = true
      }
      hostname = {
        hostname           = "https://${var.keycloak_address}"
        strict             = true
        backchannelDynamic = false
      }
      networkPolicy = {
        enabled = false
      }
    }
  }
}

resource "kubernetes_manifest" "keycloak_service_monitor" {
  manifest = {
    apiVersion = "monitoring.coreos.com/v1"
    kind       = "ServiceMonitor"
    metadata = {
      name      = "keycloak"
      namespace = var.namespace
      labels = {
        release = "prometheus"
      }
    }
    spec = {
      selector = {
        matchLabels = {
          app = "keycloak"
        }
      }
      endpoints = [
        {
          port     = "management"
          scheme   = "http"
          path     = "/metrics"
          interval = "30s"
        }
      ]
    }
  }
  depends_on = [kubernetes_manifest.keycloak]
}
