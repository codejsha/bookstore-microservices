data "kubernetes_service_account" "mysql_service_accounts" {
  for_each = toset(var.mysql_services)
  metadata {
    name      = "${each.key}-mysql"
    namespace = var.namespace
  }
}
