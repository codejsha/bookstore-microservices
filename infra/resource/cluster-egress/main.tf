terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
  }
}

locals {
  internal_nets = concat([var.pod_cidr, var.service_cidr], var.node_addresses)

  namespace_set = { for ns in var.internet_namespaces : ns => "'${ns}'" }
}

resource "kubernetes_manifest" "allow_internal" {
  manifest = {
    apiVersion = "crd.projectcalico.org/v1"
    kind       = "GlobalNetworkPolicy"
    metadata = {
      name = "cluster-egress-allow-internal"
    }
    spec = {
      order = 100
      types = ["Egress"]
      egress = [
        {
          action = "Allow"
          destination = {
            nets = local.internal_nets
          }
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "allow_exempt_namespaces" {
  manifest = {
    apiVersion = "crd.projectcalico.org/v1"
    kind       = "GlobalNetworkPolicy"
    metadata = {
      name = "cluster-egress-allow-exempt-namespaces"
    }
    spec = {
      order             = 110
      types             = ["Egress"]
      namespaceSelector = "kubernetes.io/metadata.name in {${join(", ", [for ns in var.exempt_namespaces : "'${ns}'"])}}"
      egress = [
        { action = "Allow" }
      ]
    }
  }
}

resource "kubernetes_manifest" "allow_internet_namespaces" {
  manifest = {
    apiVersion = "crd.projectcalico.org/v1"
    kind       = "GlobalNetworkPolicy"
    metadata = {
      name = "cluster-egress-allow-internet"
    }
    spec = {
      order             = 200
      types             = ["Egress"]
      namespaceSelector = "kubernetes.io/metadata.name in {${join(", ", values(local.namespace_set))}}"
      egress = [
        {
          action   = "Allow"
          protocol = "TCP"
          destination = {
            ports = var.internet_ports
          }
        },
        {
          action   = "Allow"
          protocol = "UDP"
          destination = {
            ports = [53]
          }
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "deny_default" {
  count = var.enforce ? 1 : 0

  manifest = {
    apiVersion = "crd.projectcalico.org/v1"
    kind       = "GlobalNetworkPolicy"
    metadata = {
      name = "cluster-egress-deny-default"
    }
    spec = {
      order = 1000
      types = ["Egress"]
      egress = [
        { action = "Deny" }
      ]
    }
  }

  depends_on = [
    kubernetes_manifest.allow_internal,
    kubernetes_manifest.allow_exempt_namespaces,
    kubernetes_manifest.allow_internet_namespaces,
  ]
}
