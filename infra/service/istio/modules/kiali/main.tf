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

resource "helm_release" "kiali" {
  namespace  = var.namespace
  name       = "kiali"
  repository = "https://kiali.org/helm-charts"
  chart      = "kiali-server"
  version    = var.kiali_chart_version
  values = [
    file("${path.module}/values.yaml")
  ]

  set_sensitive = [
    {
      name  = "external_services.grafana.auth.username"
      value = var.grafana_username
    },
    {
      name  = "external_services.grafana.auth.password"
      value = var.grafana_password
    }
  ]
}

locals {
  mesh_api_groups = [
    "networking.istio.io",
    "security.istio.io",
    "extensions.istio.io",
    "telemetry.istio.io",
    "gateway.networking.k8s.io",
    "inference.networking.k8s.io",
  ]
}

resource "kubernetes_cluster_role_v1" "oidc_viewer" {
  metadata {
    name = "kiali-oidc-viewer"
  }

  rule {
    api_groups = [""]
    resources  = ["configmaps", "endpoints", "pods/log"]
    verbs      = ["get", "list", "watch"]
  }

  rule {
    api_groups = [""]
    resources  = ["namespaces", "pods", "replicationcontrollers", "services"]
    verbs      = ["get", "list", "watch"]
  }

  rule {
    api_groups = ["apps"]
    resources  = ["daemonsets", "deployments", "replicasets", "statefulsets"]
    verbs      = ["get", "list", "watch"]
  }

  rule {
    api_groups = ["batch"]
    resources  = ["cronjobs", "jobs"]
    verbs      = ["get", "list", "watch"]
  }

  rule {
    api_groups = local.mesh_api_groups
    resources  = ["*"]
    verbs      = ["get", "list", "watch"]
  }

  rule {
    api_groups = ["admissionregistration.k8s.io"]
    resources  = ["mutatingwebhookconfigurations"]
    verbs      = ["get", "list", "watch"]
  }
}

resource "kubernetes_cluster_role_v1" "oidc_editor" {
  metadata {
    name = "kiali-oidc-editor"
  }

  rule {
    api_groups = [""]
    resources  = ["pods/portforward"]
    verbs      = ["create", "post"]
  }

  rule {
    api_groups = [""]
    resources  = ["namespaces", "pods", "replicationcontrollers", "services"]
    verbs      = ["patch"]
  }

  rule {
    api_groups = ["apps"]
    resources  = ["daemonsets", "deployments", "replicasets", "statefulsets"]
    verbs      = ["patch"]
  }

  rule {
    api_groups = ["batch"]
    resources  = ["cronjobs", "jobs"]
    verbs      = ["patch"]
  }

  rule {
    api_groups = local.mesh_api_groups
    resources  = ["*"]
    verbs      = ["create", "delete", "patch"]
  }
}

resource "kubernetes_cluster_role_binding_v1" "oidc_viewer" {
  metadata {
    name = "kiali-oidc-viewer"
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role_v1.oidc_viewer.metadata[0].name
  }

  dynamic "subject" {
    for_each = var.viewer_roles
    content {
      api_group = "rbac.authorization.k8s.io"
      kind      = "Group"
      name      = "${var.oidc_group_prefix}${subject.value}"
    }
  }
}

resource "kubernetes_cluster_role_binding_v1" "oidc_editor" {
  metadata {
    name = "kiali-oidc-editor"
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role_v1.oidc_editor.metadata[0].name
  }

  dynamic "subject" {
    for_each = var.editor_roles
    content {
      api_group = "rbac.authorization.k8s.io"
      kind      = "Group"
      name      = "${var.oidc_group_prefix}${subject.value}"
    }
  }
}
