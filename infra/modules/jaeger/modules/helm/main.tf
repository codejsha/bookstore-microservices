resource "helm_release" "jaeger" {
  namespace  = var.namespace
  name       = "jaeger"
  repository = "https://jaegertracing.github.io/helm-charts"
  chart      = "jaeger"
  version    = "3.4.1"
  values = [
    file("${path.module}/values-memory.yaml")
  ]
}
