resource "helm_release" "kiali" {
  namespace  = var.namespace
  name       = "kiali"
  repository = "https://kiali.org/helm-charts"
  chart      = "kiali-server"
  values = [
    file("${path.module}/values.yaml")
  ]
}
