locals {
  istio_version = "1.24.3"
}

resource "helm_release" "istio_base" {
  namespace  = var.namespace
  name       = "istio-base"
  repository = "https://istio-release.storage.googleapis.com/charts"
  chart      = "base"
  version    = local.istio_version
  set {
    name  = "defaultRevision"
    value = "default"
  }
}

resource "helm_release" "istio_istiod" {
  depends_on = [helm_release.istio_base]
  namespace  = var.namespace
  name       = "istio-istiod"
  repository = "https://istio-release.storage.googleapis.com/charts"
  chart      = "istiod"
  version    = local.istio_version
  set {
    name  = "pilot.cni.enabled"
    value = "true"
  }
  set {
    name  = "global.proxy.resources.requests.cpu"
    value = "10m"
  }
}

resource "helm_release" "istio_cni" {
  depends_on = [helm_release.istio_istiod]
  namespace  = var.namespace
  name       = "istio-cni"
  repository = "https://istio-release.storage.googleapis.com/charts"
  chart      = "cni"
  version    = local.istio_version
  # set {
  #   name  = "global.platform"
  #   value = "k3s"
  # }
}

resource "helm_release" "istio_ingressgateway" {
  depends_on = [helm_release.istio_istiod]
  namespace  = var.namespace
  name       = "istio-ingressgateway"
  repository = "https://istio-release.storage.googleapis.com/charts"
  chart      = "gateway"
  version    = local.istio_version
}
