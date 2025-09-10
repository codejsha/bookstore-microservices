terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "mysql_operator" {
  namespace  = var.namespace
  name       = "mysql-operator"
  repository = "https://mysql.github.io/mysql-operator/"
  chart      = "mysql-operator"
  version    = "2.2.8"
  values = [
    file("${path.module}/values.yaml")
  ]
}
