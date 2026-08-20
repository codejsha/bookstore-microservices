#!/usr/bin/env python3
"""Generate per-service Grafana dashboards.

Source-of-truth for the bookstore service dashboards. Run after editing this
file to regenerate every JSON under "dashboards/Bookstore - Services/".

    cd infra/service/grafana
    python3 scripts/generate_service_dashboards.py

Metric names assume Prometheus's OTLP receiver translation strategy
NoUTF8EscapingWithSuffixes (set in infra/service/prometheus/modules/helm/values.yaml).
With that strategy, OTel "http.server.request.duration" arrives as
"http_server_request_duration_seconds_{bucket,sum,count}", and resource
attributes promoted by the receiver land as labels with dots replaced by
underscores (service_name, k8s_namespace_name, k8s_deployment_name, ...).
"""

from __future__ import annotations

import copy
import json
from pathlib import Path

OUT_DIR = Path(__file__).resolve().parent.parent / "dashboards" / "Bookstore - Services"

DATASOURCE = {"type": "prometheus", "uid": "prometheus"}

SERVICES: list[dict] = [
    {"name": "catalog", "lang": "go", "grpc_server": True, "grpc_client": True, "db": "postgres"},
    {"name": "customer", "lang": "go", "grpc_server": False, "grpc_client": True, "db": "postgres"},
    {"name": "identity", "lang": "go", "grpc_server": True, "grpc_client": False, "db": "postgres"},
    {"name": "inventory", "lang": "go", "grpc_server": True, "grpc_client": False, "db": "postgres", "temporal": True},
    {"name": "admin", "lang": "jvm", "grpc_server": False, "grpc_client": False},
    {"name": "order", "lang": "jvm", "grpc_server": True, "grpc_client": True, "db": "mysql", "temporal": True},
    {"name": "payment", "lang": "jvm", "grpc_server": True, "grpc_client": False, "db": "mysql", "temporal": True},
    {"name": "settlement", "lang": "jvm", "grpc_server": False, "grpc_client": False, "db": "mysql"},
    {"name": "delivery", "lang": "python", "grpc_server": False, "grpc_client": False, "db": "mysql"},
    {"name": "notification", "lang": "python", "grpc_server": False, "grpc_client": False, "db": "mysql"},
    {"name": "support", "lang": "python", "grpc_server": True, "grpc_client": False, "db": "mysql"},
]


def panel_id() -> int:
    panel_id.counter += 1
    return panel_id.counter


panel_id.counter = 0  # type: ignore[attr-defined]


def reset_panel_ids() -> None:
    panel_id.counter = 0  # type: ignore[attr-defined]


def grid(x: int, y: int, w: int, h: int) -> dict:
    return {"x": x, "y": y, "w": w, "h": h}


def row(title: str, y: int) -> dict:
    return {
        "id": panel_id(),
        "type": "row",
        "title": title,
        "collapsed": False,
        "gridPos": grid(0, y, 24, 1),
        "panels": [],
    }


def stat_panel(title: str, expr: str, x: int, y: int, w: int = 6, h: int = 4, unit: str = "short") -> dict:
    return {
        "id": panel_id(),
        "type": "stat",
        "title": title,
        "datasource": DATASOURCE,
        "gridPos": grid(x, y, w, h),
        "fieldConfig": {
            "defaults": {
                "unit": unit,
                "color": {"mode": "thresholds"},
                "thresholds": {"mode": "absolute", "steps": [{"color": "green", "value": None}]},
            },
            "overrides": [],
        },
        "options": {
            "reduceOptions": {"calcs": ["lastNotNull"], "fields": "", "values": False},
            "textMode": "auto",
            "graphMode": "area",
            "colorMode": "value",
            "orientation": "auto",
        },
        "targets": [{"datasource": DATASOURCE, "expr": expr, "refId": "A"}],
    }


def timeseries_panel(
    title: str,
    targets: list[tuple[str, str]],
    x: int,
    y: int,
    w: int = 12,
    h: int = 8,
    unit: str = "short",
    legend: str = "list",
    stack: bool = False,
) -> dict:
    return {
        "id": panel_id(),
        "type": "timeseries",
        "title": title,
        "datasource": DATASOURCE,
        "gridPos": grid(x, y, w, h),
        "fieldConfig": {
            "defaults": {
                "unit": unit,
                "custom": {
                    "drawStyle": "line",
                    "lineInterpolation": "smooth",
                    "lineWidth": 1,
                    "fillOpacity": 10 if stack else 0,
                    "stacking": {"mode": "normal" if stack else "none", "group": "A"},
                    "axisGridShow": True,
                    "spanNulls": False,
                },
            },
            "overrides": [],
        },
        "options": {
            "legend": {"displayMode": legend, "placement": "bottom", "showLegend": True, "calcs": ["mean", "max"]},
            "tooltip": {"mode": "multi", "sort": "desc"},
        },
        "targets": [
            {"datasource": DATASOURCE, "expr": expr, "legendFormat": legend_fmt, "refId": chr(ord("A") + i)}
            for i, (expr, legend_fmt) in enumerate(targets)
        ],
    }


def red_panels(svc: str, y: int) -> list[dict]:
    """RED metrics — request rate, errors, p95/p99 latency. HTTP server."""
    sel = 'service_name=~"$service_name", k8s_namespace_name=~"$namespace"'
    panels = [
        stat_panel(
            "Request rate (rps)",
            f'sum(rate(http_server_request_duration_seconds_count{{{sel}}}[$__rate_interval]))',
            0, y, w=4, h=4, unit="reqps",
        ),
        stat_panel(
            "Error rate (%)",
            f'100 * sum(rate(http_server_request_duration_seconds_count{{{sel}, http_response_status_code=~"5..|4.."}}[$__rate_interval])) '
            f'/ sum(rate(http_server_request_duration_seconds_count{{{sel}}}[$__rate_interval]))',
            4, y, w=4, h=4, unit="percent",
        ),
        stat_panel(
            "p95 latency",
            f'histogram_quantile(0.95, sum by (le) (rate(http_server_request_duration_seconds_bucket{{{sel}}}[$__rate_interval])))',
            8, y, w=4, h=4, unit="s",
        ),
        stat_panel(
            "p99 latency",
            f'histogram_quantile(0.99, sum by (le) (rate(http_server_request_duration_seconds_bucket{{{sel}}}[$__rate_interval])))',
            12, y, w=4, h=4, unit="s",
        ),
        stat_panel(
            "Pods up",
            f'count(count by (k8s_pod_name) (http_server_request_duration_seconds_count{{{sel}}}))',
            16, y, w=4, h=4, unit="short",
        ),
        stat_panel(
            "5xx (last 5m)",
            f'sum(increase(http_server_request_duration_seconds_count{{{sel}, http_response_status_code=~"5.."}}[5m]))',
            20, y, w=4, h=4, unit="short",
        ),
    ]
    panels += [
        timeseries_panel(
            "Requests / sec by status",
            [(
                f'sum by (http_response_status_code) (rate(http_server_request_duration_seconds_count{{{sel}}}[$__rate_interval]))',
                "{{http_response_status_code}}",
            )],
            0, y + 4, w=12, h=8, unit="reqps", stack=True,
        ),
        timeseries_panel(
            "Latency percentiles",
            [
                (f'histogram_quantile(0.50, sum by (le) (rate(http_server_request_duration_seconds_bucket{{{sel}}}[$__rate_interval])))', "p50"),
                (f'histogram_quantile(0.95, sum by (le) (rate(http_server_request_duration_seconds_bucket{{{sel}}}[$__rate_interval])))', "p95"),
                (f'histogram_quantile(0.99, sum by (le) (rate(http_server_request_duration_seconds_bucket{{{sel}}}[$__rate_interval])))', "p99"),
            ],
            12, y + 4, w=12, h=8, unit="s",
        ),
    ]
    return panels


def top_routes_panels(y: int) -> list[dict]:
    sel = 'service_name=~"$service_name", k8s_namespace_name=~"$namespace"'
    return [
        timeseries_panel(
            "Top routes by rate",
            [(
                f'topk(10, sum by (http_route) (rate(http_server_request_duration_seconds_count{{{sel}}}[$__rate_interval])))',
                "{{http_route}}",
            )],
            0, y, w=12, h=8, unit="reqps",
        ),
        timeseries_panel(
            "p95 latency by route",
            [(
                f'topk(10, histogram_quantile(0.95, sum by (le, http_route) '
                f'(rate(http_server_request_duration_seconds_bucket{{{sel}}}[$__rate_interval]))))',
                "{{http_route}}",
            )],
            12, y, w=12, h=8, unit="s",
        ),
    ]


def grpc_server_panels(y: int) -> list[dict]:
    sel = 'service_name=~"$service_name", k8s_namespace_name=~"$namespace"'
    return [
        timeseries_panel(
            "gRPC server: requests / sec by method",
            [(
                f'sum by (rpc_service, rpc_method) (rate(rpc_server_duration_milliseconds_count{{{sel}}}[$__rate_interval]))',
                "{{rpc_service}}/{{rpc_method}}",
            )],
            0, y, w=12, h=8, unit="reqps",
        ),
        timeseries_panel(
            "gRPC server: p95 latency by method",
            [(
                f'histogram_quantile(0.95, sum by (le, rpc_service, rpc_method) '
                f'(rate(rpc_server_duration_milliseconds_bucket{{{sel}}}[$__rate_interval])))',
                "{{rpc_service}}/{{rpc_method}}",
            )],
            12, y, w=12, h=8, unit="ms",
        ),
        timeseries_panel(
            "gRPC server: non-OK responses",
            [(
                f'sum by (rpc_grpc_status_code) (rate(rpc_server_duration_milliseconds_count{{{sel}, rpc_grpc_status_code!="0"}}[$__rate_interval]))',
                "{{rpc_grpc_status_code}}",
            )],
            0, y + 8, w=24, h=6, unit="reqps", stack=True,
        ),
    ]


def grpc_client_panels(y: int) -> list[dict]:
    sel = 'service_name=~"$service_name", k8s_namespace_name=~"$namespace"'
    return [
        timeseries_panel(
            "gRPC client: requests / sec by peer",
            [(
                f'sum by (rpc_service, rpc_method) (rate(rpc_client_duration_milliseconds_count{{{sel}}}[$__rate_interval]))',
                "{{rpc_service}}/{{rpc_method}}",
            )],
            0, y, w=12, h=8, unit="reqps",
        ),
        timeseries_panel(
            "gRPC client: p95 latency by peer",
            [(
                f'histogram_quantile(0.95, sum by (le, rpc_service, rpc_method) '
                f'(rate(rpc_client_duration_milliseconds_bucket{{{sel}}}[$__rate_interval])))',
                "{{rpc_service}}/{{rpc_method}}",
            )],
            12, y, w=12, h=8, unit="ms",
        ),
    ]


def jvm_panels(y: int) -> list[dict]:
    sel = 'service_name=~"$service_name", k8s_namespace_name=~"$namespace"'
    return [
        timeseries_panel(
            "JVM heap by area",
            [(
                f'sum by (jvm_memory_pool_name) (jvm_memory_used_bytes{{{sel}, jvm_memory_type="heap"}})',
                "{{jvm_memory_pool_name}}",
            )],
            0, y, w=12, h=8, unit="bytes", stack=True,
        ),
        timeseries_panel(
            "JVM non-heap by area",
            [(
                f'sum by (jvm_memory_pool_name) (jvm_memory_used_bytes{{{sel}, jvm_memory_type="non_heap"}})',
                "{{jvm_memory_pool_name}}",
            )],
            12, y, w=12, h=8, unit="bytes", stack=True,
        ),
        timeseries_panel(
            "GC pause time",
            [(
                f'sum by (jvm_gc_action, jvm_gc_name) (rate(jvm_gc_duration_seconds_sum{{{sel}}}[$__rate_interval]))',
                "{{jvm_gc_name}} ({{jvm_gc_action}})",
            )],
            0, y + 8, w=12, h=6, unit="s",
        ),
        timeseries_panel(
            "Threads & classes",
            [
                (f'sum by (jvm_thread_state) (jvm_thread_count{{{sel}}})', "threads {{jvm_thread_state}}"),
                (f'sum(jvm_class_loaded_count{{{sel}}})', "classes loaded"),
            ],
            12, y + 8, w=12, h=6, unit="short",
        ),
    ]


def hikari_panels(y: int) -> list[dict]:
    sel = 'service_name=~"$service_name", k8s_namespace_name=~"$namespace"'
    return [
        timeseries_panel(
            "HikariCP connections",
            [
                (f'sum by (pool_name) (hikaricp_connections_active{{{sel}}})', "{{pool_name}} active"),
                (f'sum by (pool_name) (hikaricp_connections_idle{{{sel}}})', "{{pool_name}} idle"),
                (f'sum by (pool_name) (hikaricp_connections_pending{{{sel}}})', "{{pool_name}} pending"),
            ],
            0, y, w=12, h=8, unit="short",
        ),
        timeseries_panel(
            "Connection acquire time (p95)",
            [(
                f'histogram_quantile(0.95, sum by (le, pool_name) '
                f'(rate(hikaricp_connections_acquire_seconds_bucket{{{sel}}}[$__rate_interval])))',
                "{{pool_name}}",
            )],
            12, y, w=12, h=8, unit="s",
        ),
    ]


def go_runtime_panels(y: int) -> list[dict]:
    sel = 'service_name=~"$service_name", k8s_namespace_name=~"$namespace"'
    return [
        timeseries_panel(
            "Goroutines",
            [(f'sum by (k8s_pod_name) (process_runtime_go_goroutines{{{sel}}})', "{{k8s_pod_name}}")],
            0, y, w=8, h=6, unit="short",
        ),
        timeseries_panel(
            "Heap allocated (live)",
            [(f'sum by (k8s_pod_name) (process_runtime_go_mem_heap_alloc_bytes{{{sel}}})', "{{k8s_pod_name}}")],
            8, y, w=8, h=6, unit="bytes",
        ),
        timeseries_panel(
            "GC pause / sec",
            [(f'sum by (k8s_pod_name) (rate(process_runtime_go_gc_pause_ns_sum{{{sel}}}[$__rate_interval])) / 1e9', "{{k8s_pod_name}}")],
            16, y, w=8, h=6, unit="s",
        ),
    ]


def python_runtime_panels(y: int) -> list[dict]:
    sel = 'service_name=~"$service_name", k8s_namespace_name=~"$namespace"'
    return [
        timeseries_panel(
            "Active requests (in flight)",
            [(f'sum by (k8s_pod_name) (http_server_active_requests{{{sel}}})', "{{k8s_pod_name}}")],
            0, y, w=12, h=6, unit="short",
        ),
        timeseries_panel(
            "Process CPU time / sec",
            [(f'sum by (k8s_pod_name) (rate(process_cpu_time_seconds_total{{{sel}}}[$__rate_interval]))', "{{k8s_pod_name}}")],
            12, y, w=12, h=6, unit="percentunit",
        ),
    ]


def temporal_panels(y: int) -> list[dict]:
    """Temporal SDK metrics emitted by workers."""
    sel = 'service_name=~"$service_name", k8s_namespace_name=~"$namespace"'
    return [
        timeseries_panel(
            "Workflow task latency (p95)",
            [(
                f'histogram_quantile(0.95, sum by (le, workflow_type) '
                f'(rate(temporal_workflow_task_execution_latency_bucket{{{sel}}}[$__rate_interval])))',
                "{{workflow_type}}",
            )],
            0, y, w=12, h=8, unit="s",
        ),
        timeseries_panel(
            "Activity execution latency (p95)",
            [(
                f'histogram_quantile(0.95, sum by (le, activity_type) '
                f'(rate(temporal_activity_execution_latency_bucket{{{sel}}}[$__rate_interval])))',
                "{{activity_type}}",
            )],
            12, y, w=12, h=8, unit="s",
        ),
        timeseries_panel(
            "Workflow / activity failures",
            [
                (f'sum by (workflow_type) (rate(temporal_workflow_failed{{{sel}}}[$__rate_interval]))', "wf {{workflow_type}}"),
                (f'sum by (activity_type) (rate(temporal_activity_execution_failed{{{sel}}}[$__rate_interval]))', "act {{activity_type}}"),
            ],
            0, y + 8, w=24, h=6, unit="short",
        ),
    ]


def pod_health_panels(y: int) -> list[dict]:
    sel = 'k8s_namespace_name=~"$namespace", k8s_deployment_name=~"$service_name"'
    return [
        timeseries_panel(
            "Container CPU",
            [(
                f'sum by (k8s_pod_name) (rate(container_cpu_usage_seconds_total{{namespace=~"$namespace", pod=~".*$service_name.*"}}[$__rate_interval]))',
                "{{k8s_pod_name}}",
            )],
            0, y, w=12, h=6, unit="percentunit",
        ),
        timeseries_panel(
            "Container memory (RSS)",
            [(
                f'sum by (k8s_pod_name) (container_memory_working_set_bytes{{namespace=~"$namespace", pod=~".*$service_name.*", image!=""}})',
                "{{k8s_pod_name}}",
            )],
            12, y, w=12, h=6, unit="bytes",
        ),
        timeseries_panel(
            "Pod restarts (1h)",
            [(
                f'sum by (pod) (increase(kube_pod_container_status_restarts_total{{namespace=~"$namespace", pod=~".*$service_name.*"}}[1h]))',
                "{{pod}}",
            )],
            0, y + 6, w=24, h=4, unit="short",
        ),
    ]


def template_vars(default_service: str) -> list[dict]:
    return [
        {
            "name": "namespace",
            "label": "Namespace",
            "type": "query",
            "datasource": DATASOURCE,
            "query": "label_values(http_server_request_duration_seconds_count, k8s_namespace_name)",
            "current": {"text": "All", "value": "$__all"},
            "includeAll": True,
            "multi": True,
            "refresh": 2,
            "sort": 1,
        },
        {
            "name": "service_name",
            "label": "Service",
            "type": "query",
            "datasource": DATASOURCE,
            "query": "label_values(http_server_request_duration_seconds_count{k8s_namespace_name=~\"$namespace\"}, service_name)",
            "current": {"text": default_service, "value": default_service},
            "includeAll": False,
            "multi": False,
            "refresh": 2,
            "sort": 1,
        },
    ]


def build_dashboard(svc: dict) -> dict:
    reset_panel_ids()
    name = svc["name"]
    panels: list[dict] = []
    y = 0

    panels.append(row("Golden signals (HTTP)", y)); y += 1
    panels += red_panels(name, y); y += 12

    panels.append(row("Routes", y)); y += 1
    panels += top_routes_panels(y); y += 8

    if svc.get("grpc_server"):
        panels.append(row("gRPC server", y)); y += 1
        panels += grpc_server_panels(y); y += 14

    if svc.get("grpc_client"):
        panels.append(row("gRPC client", y)); y += 1
        panels += grpc_client_panels(y); y += 8

    if svc["lang"] == "jvm":
        panels.append(row("JVM", y)); y += 1
        panels += jvm_panels(y); y += 14
        if svc.get("db"):
            panels.append(row("Database (HikariCP)", y)); y += 1
            panels += hikari_panels(y); y += 8
    elif svc["lang"] == "go":
        panels.append(row("Go runtime", y)); y += 1
        panels += go_runtime_panels(y); y += 6
    elif svc["lang"] == "python":
        panels.append(row("Python runtime", y)); y += 1
        panels += python_runtime_panels(y); y += 6

    if svc.get("temporal"):
        panels.append(row("Temporal", y)); y += 1
        panels += temporal_panels(y); y += 14

    panels.append(row("Pod health", y)); y += 1
    panels += pod_health_panels(y); y += 10

    return {
        "title": f"Bookstore — {name.title()}",
        "uid": f"bookstore-{name}",
        "tags": ["bookstore", "service", name],
        "schemaVersion": 39,
        "version": 1,
        "editable": True,
        "graphTooltip": 1,
        "time": {"from": "now-6h", "to": "now"},
        "timepicker": {},
        "refresh": "30s",
        "annotations": {"list": []},
        "templating": {"list": template_vars(name)},
        "panels": panels,
    }


def build_overview() -> dict:
    """Fleet-wide RED dashboard — repeats one row per service."""
    reset_panel_ids()
    sel = 'k8s_namespace_name=~"$namespace"'
    panels: list[dict] = []

    services_filter = "|".join(s["name"] for s in SERVICES)

    panels.append({
        "id": panel_id(),
        "type": "row",
        "title": "Fleet RED",
        "collapsed": False,
        "gridPos": grid(0, 0, 24, 1),
        "panels": [],
    })

    panels.append(timeseries_panel(
        "Requests / sec by service",
        [(
            f'sum by (service_name) (rate(http_server_request_duration_seconds_count{{{sel}, service_name=~"{services_filter}"}}[$__rate_interval]))',
            "{{service_name}}",
        )],
        0, 1, w=12, h=8, unit="reqps", stack=True,
    ))

    panels.append(timeseries_panel(
        "Error rate (%) by service",
        [(
            f'100 * sum by (service_name) (rate(http_server_request_duration_seconds_count{{{sel}, service_name=~"{services_filter}", http_response_status_code=~"5.."}}[$__rate_interval])) '
            f'/ sum by (service_name) (rate(http_server_request_duration_seconds_count{{{sel}, service_name=~"{services_filter}"}}[$__rate_interval]))',
            "{{service_name}}",
        )],
        12, 1, w=12, h=8, unit="percent",
    ))

    panels.append(timeseries_panel(
        "p95 latency by service",
        [(
            f'histogram_quantile(0.95, sum by (le, service_name) (rate(http_server_request_duration_seconds_bucket{{{sel}, service_name=~"{services_filter}"}}[$__rate_interval])))',
            "{{service_name}}",
        )],
        0, 9, w=12, h=8, unit="s",
    ))

    panels.append(timeseries_panel(
        "p99 latency by service",
        [(
            f'histogram_quantile(0.99, sum by (le, service_name) (rate(http_server_request_duration_seconds_bucket{{{sel}, service_name=~"{services_filter}"}}[$__rate_interval])))',
            "{{service_name}}",
        )],
        12, 9, w=12, h=8, unit="s",
    ))

    panels.append({
        "id": panel_id(),
        "type": "row",
        "title": "gRPC fleet",
        "collapsed": False,
        "gridPos": grid(0, 17, 24, 1),
        "panels": [],
    })

    panels.append(timeseries_panel(
        "gRPC server requests / sec",
        [(
            f'sum by (service_name) (rate(rpc_server_duration_milliseconds_count{{{sel}, service_name=~"{services_filter}"}}[$__rate_interval]))',
            "{{service_name}}",
        )],
        0, 18, w=12, h=8, unit="reqps", stack=True,
    ))

    panels.append(timeseries_panel(
        "gRPC server non-OK / sec",
        [(
            f'sum by (service_name) (rate(rpc_server_duration_milliseconds_count{{{sel}, service_name=~"{services_filter}", rpc_grpc_status_code!="0"}}[$__rate_interval]))',
            "{{service_name}}",
        )],
        12, 18, w=12, h=8, unit="reqps",
    ))

    return {
        "title": "Bookstore — Fleet overview",
        "uid": "bookstore-overview",
        "tags": ["bookstore", "overview"],
        "schemaVersion": 39,
        "version": 1,
        "editable": True,
        "graphTooltip": 1,
        "time": {"from": "now-6h", "to": "now"},
        "refresh": "30s",
        "annotations": {"list": []},
        "templating": {
            "list": [
                {
                    "name": "namespace",
                    "label": "Namespace",
                    "type": "query",
                    "datasource": DATASOURCE,
                    "query": "label_values(http_server_request_duration_seconds_count, k8s_namespace_name)",
                    "current": {"text": "All", "value": "$__all"},
                    "includeAll": True,
                    "multi": True,
                    "refresh": 2,
                    "sort": 1,
                }
            ]
        },
        "panels": panels,
    }


def main() -> None:
    OUT_DIR.mkdir(parents=True, exist_ok=True)
    overview = build_overview()
    (OUT_DIR / "overview.json").write_text(json.dumps(overview, indent=2) + "\n")
    for svc in SERVICES:
        dash = build_dashboard(copy.deepcopy(svc))
        (OUT_DIR / f"{svc['name']}.json").write_text(json.dumps(dash, indent=2) + "\n")
    print(f"Wrote {1 + len(SERVICES)} dashboards into {OUT_DIR}")


if __name__ == "__main__":
    main()
