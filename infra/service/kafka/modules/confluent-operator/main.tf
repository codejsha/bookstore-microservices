terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

resource "kubernetes_service_account_v1" "confluent_operator" {
  metadata {
    name      = "confluent-operator"
    namespace = var.namespace
  }
}

resource "kubernetes_secret_v1" "confluent_operator" {
  type = "kubernetes.io/service-account-token"
  metadata {
    name      = "confluent-operator-token"
    namespace = var.namespace
    annotations = {
      "kubernetes.io/service-account.name" = kubernetes_service_account_v1.confluent_operator.metadata[0].name
    }
  }
}

resource "helm_release" "confluent_operator" {
  namespace  = var.namespace
  name       = "confluent-operator"
  repository = "https://packages.confluent.io/helm"
  chart      = "confluent-for-kubernetes"
  version    = "0.1514.1"
  values = [
    file("${path.module}/values.yaml")
  ]
  set_list = [
    { name = "namespaceList", value = var.watch_namespaces }
  ]
}
