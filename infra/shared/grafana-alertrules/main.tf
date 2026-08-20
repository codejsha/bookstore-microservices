terraform {
  required_providers {
    grafana = {
      source = "grafana/grafana"
    }
  }
}

resource "grafana_folder" "this" {
  title             = var.folder_title
  parent_folder_uid = var.parent_folder_uid
}

resource "grafana_rule_group" "this" {
  name             = var.group_name
  folder_uid       = grafana_folder.this.uid
  interval_seconds = var.interval_seconds

  dynamic "rule" {
    for_each = var.rules
    content {
      name           = rule.value.name
      condition      = "C"
      for            = rule.value.for
      no_data_state  = rule.value.no_data_state
      exec_err_state = "Error"
      labels         = { severity = rule.value.severity }
      annotations = {
        summary     = rule.value.summary
        description = rule.value.description
      }

      data {
        ref_id         = "A"
        datasource_uid = coalesce(rule.value.datasource_uid, var.datasource_uid)
        relative_time_range {
          from = 600
          to   = 0
        }
        model = jsonencode(merge(
          {
            refId      = "A"
            datasource = { type = coalesce(rule.value.datasource_type, var.datasource_type), uid = coalesce(rule.value.datasource_uid, var.datasource_uid) }
            expr       = rule.value.expr
          },
          coalesce(rule.value.datasource_type, var.datasource_type) == "loki"
          ? { queryType = "instant", editorMode = "code" }
          : { instant = true, range = false, intervalMs = 1000, maxDataPoints = 43200 }
        ))
      }

      data {
        ref_id         = "B"
        datasource_uid = "__expr__"
        relative_time_range {
          from = 600
          to   = 0
        }
        model = jsonencode({
          refId      = "B"
          type       = "reduce"
          datasource = { type = "__expr__", uid = "__expr__" }
          expression = "A"
          reducer    = "last"
        })
      }

      data {
        ref_id         = "C"
        datasource_uid = "__expr__"
        relative_time_range {
          from = 600
          to   = 0
        }
        model = jsonencode({
          refId      = "C"
          type       = "threshold"
          datasource = { type = "__expr__", uid = "__expr__" }
          expression = "B"
          conditions = [{ evaluator = { type = "gt", params = [0] } }]
        })
      }
    }
  }
}
