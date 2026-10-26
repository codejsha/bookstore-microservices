include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "kubernetes" {
  path = find_in_parent_folders("_envcommon/provider_kubernetes.hcl")
}

include "helm" {
  path = find_in_parent_folders("_envcommon/provider_helm.hcl")
}

inputs = {
  namespace = "cloud-config"
  repo_root = get_repo_root()

  client_principals = [
    "cluster.local/ns/bookstore/sa/admin",
    "cluster.local/ns/bookstore/sa/catalog-postgres",
    "cluster.local/ns/bookstore/sa/customer-postgres",
    "cluster.local/ns/bookstore/sa/identity-postgres",
    "cluster.local/ns/bookstore/sa/fulfillment-postgres",
    "cluster.local/ns/bookstore/sa/order-postgres",
    "cluster.local/ns/bookstore/sa/payment-postgres",
    "cluster.local/ns/bookstore/sa/settlement-mysql",
    "cluster.local/ns/istio-system/sa/bookstore-gateway-istio",
  ]
}

terraform {
  source = "${get_repo_root()}/infra//service/cloud-config/terraform"
}
