resource "kubernetes_service_account_v1" "service_account" {
  metadata {
    name      = "${var.service_account}-sa"
    namespace = var.namespace
  }
  dynamic "image_pull_secret" {
    for_each = var.image_pull_secrets
    content {
      name = image_pull_secret.value
    }
  }
}

resource "kubernetes_role_binding_v1" "binding" {
  metadata {
    name      = "${var.service_account}-binding"
    namespace = var.namespace
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = "tekton-triggers-eventlistener-roles"
  }
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account_v1.service_account.metadata[0].name
    namespace = var.namespace
  }
}

resource "kubernetes_cluster_role_binding_v1" "cluster_binding" {
  metadata {
    name = "${var.service_account}-clusterbinding"
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = "tekton-triggers-eventlistener-clusterroles"
  }
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account_v1.service_account.metadata[0].name
    namespace = var.namespace
  }
}
