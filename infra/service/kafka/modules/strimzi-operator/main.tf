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

resource "kubernetes_service_account_v1" "strimzi_cluster_operator" {
  metadata {
    name      = "strimzi-cluster-operator"
    namespace = var.namespace
  }
}

resource "kubernetes_secret_v1" "strimzi_cluster_operator" {
  type = "kubernetes.io/service-account-token"
  metadata {
    name      = "strimzi-cluster-operator-token"
    namespace = var.namespace
    annotations = {
      "kubernetes.io/service-account.name" = kubernetes_service_account_v1.strimzi_cluster_operator.metadata[0].name
    }
  }
}

resource "helm_release" "strimzi_cluster_operator" {
  namespace = var.namespace
  name      = "strimzi-cluster-operator"
  chart     = "oci://quay.io/strimzi-helm/strimzi-kafka-operator"
  version   = "0.51.0"
  values = [
    file("${path.module}/values.yaml")
  ]
  set_list = [
    { name = "watchNamespaces", value = var.watch_namespaces }
  ]
}
