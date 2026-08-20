resource "kubernetes_service_account_v1" "mysql" {
  for_each = toset(var.mysql_services)
  metadata {
    name      = "${each.key}-mysql"
    namespace = var.namespace
  }
}

resource "kubernetes_secret_v1" "mysql" {
  for_each = toset(var.mysql_services)
  type     = "kubernetes.io/service-account-token"
  metadata {
    name      = "${each.key}-mysql-token"
    namespace = var.namespace
    annotations = {
      "kubernetes.io/service-account.name" = kubernetes_service_account_v1.mysql[each.key].metadata[0].name
    }
  }
}
