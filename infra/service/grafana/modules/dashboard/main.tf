terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

locals {
  dashboard_root  = "${path.root}/dashboards"
  dashboard_files = fileset(local.dashboard_root, "**/*.json")
  dashboards = {
    for f in local.dashboard_files :
    f => {
      content = file("${local.dashboard_root}/${f}")
      folder  = dirname(f)
      name    = trimsuffix(basename(f), ".json")
    }
  }
}

resource "kubernetes_config_map_v1" "grafana_dashboards" {
  for_each = local.dashboards

  metadata {
    name      = "grafana-dashboard-${lower(replace(replace(replace(each.key, "/", "-"), " ", "-"), ".json", ""))}"
    namespace = var.namespace
    labels = {
      grafana_dashboard = "bookstore"
    }
    annotations = each.value.folder == "." ? {} : {
      grafana_folder = each.value.folder
    }
  }

  data = {
    "${each.value.name}.json" = each.value.content
  }
}
