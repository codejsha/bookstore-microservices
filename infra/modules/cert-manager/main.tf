provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

resource "kubernetes_namespace" "cert_manager" {
  metadata {
    name = var.namespace
  }
}

resource "helm_release" "cert_manager" {
  namespace  = kubernetes_namespace.cert_manager.metadata.0.name
  name       = "cert-manager"
  repository = "https://charts.jetstack.io"
  chart      = "cert-manager"
  version    = "v1.17.0"
  set {
    name  = "crds.enabled"
    value = true
  }
}

resource "kubernetes_manifest" "cluster_issuer" {
  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "ClusterIssuer"
    metadata = {
      name = "selfsigned-cluster-issuer"
    }
    spec = {
      selfSigned = {}
    }
  }
}
