include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "kubernetes" {
  path = find_in_parent_folders("_envcommon/provider_kubernetes.hcl")
}

dependencies {
  paths = [
    "../../service/calico",
  ]
}

inputs = {
  pod_cidr       = "10.42.0.0/16"
  service_cidr   = "10.43.0.0/16"
  node_addresses = ["192.168.0.20/32"]

  exempt_namespaces   = ["kube-system"]
  internet_namespaces = ["nexus", "kafka", "bookstore-ci"]
  internet_ports      = [80, 443]

  enforce = false
}

terraform {
  source = "${get_repo_root()}/infra//resource/cluster-egress"
}
