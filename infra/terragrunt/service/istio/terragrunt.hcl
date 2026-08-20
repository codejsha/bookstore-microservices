include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "kubernetes" {
  path = find_in_parent_folders("_envcommon/provider_kubernetes.hcl")
}

include "helm" {
  path = find_in_parent_folders("_envcommon/provider_helm.hcl")
}

remote_state {
  backend = "local"

  generate = {
    path      = "backend.tf"
    if_exists = "overwrite_terragrunt"
  }

  config = {
    path = "${get_terragrunt_dir()}/terraform.tfstate"
  }
}

inputs = {
  namespace           = "istio-system"
  istio_version       = "1.29.1"
  kiali_chart_version = "2.26.0"
  kiali_address       = "kiali.example.com"
  kiali_service_name  = "kiali"
  grafana_username    = "admin"
  grafana_password    = "secret:grafana-admin-credentials:admin_password"

  edge_hosts = {
    "alloy-grpc.example.com"          = ["grafana"]
    "alloy-http.example.com"          = ["grafana"]
    "argocd.example.com"              = ["argocd"]
    "bookstore.example.com"           = ["bookstore"]
    "bookstore-admin.example.com"     = ["bookstore"]
    "bookstore-admin-api.example.com" = ["bookstore"]
    "catalog-api.example.com"         = ["bookstore"]
    "cloud-config.example.com"        = ["cloud-config", "bookstore"]
    "customer-api.example.com"        = ["bookstore"]
    "delivery-api.example.com"        = ["bookstore"]
    "git.example.com"                 = ["gitea"]
    "grafana.example.com"             = ["grafana"]
    "harbor.example.com"              = ["harbor"]
    "hyperswitch.example.com"         = ["hyperswitch"]
    "hyperswitch-api.example.com"     = ["hyperswitch"]
    "identity-api.example.com"        = ["bookstore"]
    "inventory-api.example.com"       = ["bookstore"]
    "keycloak.example.com"            = ["keycloak"]
    "kiali.example.com"               = ["istio-system"]
    "nexus.example.com"               = ["nexus"]
    "notification-api.example.com"    = ["bookstore"]
    "opensearch.example.com"          = ["opensearch"]
    "opensearch-api.example.com"      = ["opensearch"]
    "order-api.example.com"           = ["bookstore"]
    "payment-api.example.com"         = ["bookstore"]
    "prometheus.example.com"          = ["prometheus"]
    "seaweedfs.example.com"           = ["seaweedfs"]
    "seaweedfs-api.example.com"       = ["seaweedfs"]
    "settlement-api.example.com"      = ["bookstore"]
    "support-api.example.com"         = ["bookstore"]
    "tekton.example.com"              = ["tekton-pipelines"]
    "temporal.example.com"            = ["temporal"]
    "vault.example.com"               = ["vault"]
  }
}

terraform {
  source = "${get_repo_root()}/infra//service/istio"
}
