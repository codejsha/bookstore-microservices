terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

locals {
  agent_config = templatefile("${path.module}/agent.alloy.tftpl", {
    gateway_endpoint    = "${var.gateway_target_fqdn}:4317"
    pyroscope_write_url = var.pyroscope_write_url
  })

  gateway_config = templatefile("${path.module}/gateway.alloy.tftpl", {
    loki_push_url       = var.loki_push_url
    tempo_otlp_endpoint = var.tempo_otlp_endpoint
    prometheus_otlp_url = var.prometheus_otlp_url
  })
}

resource "helm_release" "alloy_agent" {
  namespace  = var.namespace
  name       = "alloy-agent"
  repository = "https://grafana.github.io/helm-charts"
  chart      = "alloy"
  version    = "0.12.6"
  values = [
    templatefile("${path.module}/agent.values.yaml", {
      config_contents = local.agent_config
    })
  ]
}

resource "helm_release" "alloy_gateway" {
  namespace  = var.namespace
  name       = "alloy-gateway"
  repository = "https://grafana.github.io/helm-charts"
  chart      = "alloy"
  version    = "0.12.6"
  values = [
    templatefile("${path.module}/gateway.values.yaml", {
      config_contents = local.gateway_config
    })
  ]
}
