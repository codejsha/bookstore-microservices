resource "kubernetes_service_account_v1" "postgres" {
  for_each = toset(var.postgres_services)
  metadata {
    name      = "${each.key}-postgres"
    namespace = var.namespace
  }
}

resource "kubernetes_secret_v1" "postgres" {
  for_each = toset(var.postgres_services)
  type     = "kubernetes.io/service-account-token"
  metadata {
    name      = "${each.key}-postgres-token"
    namespace = var.namespace
    annotations = {
      "kubernetes.io/service-account.name" = kubernetes_service_account_v1.postgres[each.key].metadata[0].name
    }
  }
}
