terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
  }
}

data "kubernetes_service_v1" "gateway" {
  metadata {
    name      = var.gateway_service_name
    namespace = var.gateway_namespace
  }
}

locals {
  gateway_ip = (
    var.gateway_address == "cluster_ip"
    ? data.kubernetes_service_v1.gateway.spec[0].cluster_ip
    : data.kubernetes_service_v1.gateway.status[0].load_balancer[0].ingress[0].ip
  )

  zone_block = <<-EOT
    ${var.edge_zone}:53 {
        errors
        cache 30
        template IN A {
            answer "{{ .Name }} ${var.record_ttl} IN A ${local.gateway_ip}"
        }
        template IN AAAA {
            rcode NOERROR
        }
        template ANY ANY {
            rcode NOERROR
        }
    }
  EOT
}

resource "kubernetes_config_map_v1" "coredns_custom" {
  metadata {
    name      = "coredns-custom"
    namespace = "kube-system"
  }
  data = {
    "${replace(var.edge_zone, ".", "-")}.server" = local.zone_block
  }
}
