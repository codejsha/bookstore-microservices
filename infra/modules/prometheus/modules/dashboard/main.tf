locals {
  dashboards = {
    "istio-mesh.json" = file("dashboards/istio-mesh.json")
    "istio-service.json" = file("dashboards/istio-service.json")
    "istio-workload.json" = file("dashboards/istio-workload.json")
    "istio-performance.json" = file("dashboards/istio-performance.json")
    "istio-controlplane.json" = file("dashboards/istio-controlplane.json")
    "istio-wasm.json" = file("dashboards/istio-wasm.json")
  }
}

resource "kubernetes_config_map" "grafana_dashboards" {
  for_each = local.dashboards

  metadata {
    name      = "grafana-dashboard-${replace(each.key, ".json", "")}"
    namespace = var.namespace
    labels = {
      grafana_dashboard = "1"
    }
  }
  data = {
    (each.key) = each.value
  }
}
