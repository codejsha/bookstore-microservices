terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

resource "helm_release" "istio_base" {
  upgrade_install = true
  namespace       = var.namespace
  name            = "istio-base"
  repository      = "https://istio-release.storage.googleapis.com/charts"
  chart           = "base"
  version         = var.istio_version
}

resource "helm_release" "istio_istiod" {
  upgrade_install = true
  depends_on      = [helm_release.istio_base]
  namespace       = var.namespace
  name            = "istio-istiod"
  repository      = "https://istio-release.storage.googleapis.com/charts"
  chart           = "istiod"
  version         = var.istio_version
  values = [
    file("${path.module}/values.yaml")
  ]
  set = [
    { name = "profile", value = "ambient" }
  ]
}

resource "helm_release" "istio_cni" {
  upgrade_install = true
  depends_on      = [helm_release.istio_istiod]
  namespace       = var.namespace
  name            = "istio-cni"
  repository      = "https://istio-release.storage.googleapis.com/charts"
  chart           = "cni"
  version         = var.istio_version
  set = [
    { name = "profile", value = "ambient" },
  ]
}

resource "helm_release" "istio_ztunnel" {
  upgrade_install = true
  depends_on      = [helm_release.istio_cni]
  namespace       = var.namespace
  name            = "istio-ztunnel"
  repository      = "https://istio-release.storage.googleapis.com/charts"
  chart           = "ztunnel"
  version         = var.istio_version
}
