resource "kubernetes_manifest" "keycloak_db" {
  manifest = {
    apiVersion = "postgresql.cnpg.io/v1"
    kind       = "Cluster"
    metadata = {
      name      = "keycloak-db"
      namespace = var.namespace
    }
    spec = {
      instances             = 1
      imageName             = "ghcr.io/cloudnative-pg/postgresql:18"
      primaryUpdateStrategy = "unsupervised"
      enableSuperuserAccess = true
      resources = {
        requests = { cpu = "50m", memory = "256Mi" }
        limits   = { cpu = "500m", memory = "1Gi" }
      }
      storage = {
        size = "5Gi"
      }
    }
  }
}
