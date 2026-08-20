terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
    http = {
      source  = "hashicorp/http"
      version = "~> 3.6"
    }
  }
}

resource "kubernetes_manifest" "cluster_issuer" {
  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "ClusterIssuer"
    metadata = {
      name = "selfsigned-cluster-issuer"
    }
    spec = {
      selfSigned = {}
    }
  }
}

data "http" "pki_ca_chain" {
  url = "${var.vault_url}/v1/${var.vault_pki_mount}/ca_chain"
}

resource "kubernetes_config_map_v1" "pki_ca_chain" {
  metadata {
    name      = "vault-pki-ca"
    namespace = var.trust_namespace
  }
  data = {
    "ca-chain.pem" = data.http.pki_ca_chain.response_body
  }
}

locals {
  ca_bundle_sources = [
    { useDefaultCAs = true },
    {
      configMap = {
        name = kubernetes_config_map_v1.pki_ca_chain.metadata[0].name
        key  = "ca-chain.pem"
      }
    },
  ]
}

resource "kubernetes_manifest" "internal_ca_bundle" {
  manifest = {
    apiVersion = "trust.cert-manager.io/v1alpha1"
    kind       = "Bundle"
    metadata = {
      name = var.bundle_name
    }
    spec = {
      sources = local.ca_bundle_sources
      target = {
        configMap = {
          key = var.bundle_key
        }
        namespaceSelector = {
          matchLabels = {}
        }
      }
    }
  }
}

resource "kubernetes_manifest" "tekton_trusted_ca_bundle" {
  manifest = {
    apiVersion = "trust.cert-manager.io/v1alpha1"
    kind       = "Bundle"
    metadata = {
      name = var.tekton_bundle_name
    }
    spec = {
      sources = local.ca_bundle_sources
      target = {
        configMap = {
          key = var.bundle_key
        }
        namespaceSelector = {
          matchLabels = {}
        }
      }
    }
  }
}
