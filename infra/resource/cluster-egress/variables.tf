variable "pod_cidr" {
  description = <<-EOT
    CIDR the Calico IPPool hands pod addresses from. calico-ipam ignores
    node.spec.podCIDR and allocates from its own default-ipv4-ippool, so the
    value must match the live IPPool (`kubectl get ippools.crd.projectcalico.org`),
    not the k3s --cluster-cidr. The default is the pinned CALICO_IPV4POOL_CIDR
    a fresh bootstrap creates; a cluster bootstrapped before that pin carries
    whatever calico-node picked and overrides this input.
  EOT
  type        = string
  default     = "10.42.0.0/16"
}

variable "service_cidr" {
  description = "k3s ClusterIP range. kube-dns (10.43.0.10) and every in-cluster Service live here"
  type        = string
  default     = "10.43.0.0/16"
}

variable "node_addresses" {
  description = <<-EOT
    Node addresses pods may still reach, as CIDRs. Needed for the kube API
    (:6443), hostNetwork/hostPort pods, and the edge gateway: cluster-dns
    answers every example.com name with the gateway's LoadBalancer IP, which
    on k3s ServiceLB is the node IP itself.
  EOT
  type        = list(string)
  default     = ["192.168.0.20/32"]
}

variable "exempt_namespaces" {
  description = <<-EOT
    Namespaces the egress policy never restricts. kube-system holds CoreDNS,
    whose upstream forward (`forward . /etc/resolv.conf`) is the only pod-side
    DNS resolution that leaves the cluster.
  EOT
  type        = list(string)
  default     = ["kube-system"]
}

variable "internet_namespaces" {
  description = <<-EOT
    Namespaces allowed to open TCP connections to any destination. Keep this
    to the workloads that fetch from the internet from inside a pod netns:
    Nexus proxy repositories, Strimzi KafkaConnect builds (repo1.maven.org,
    hub-downloads.confluent.io), and the CI docker-in-docker sidecar pulling
    Testcontainers images. Image pulls by kubelet happen in the host netns
    and are not subject to this policy.
  EOT
  type        = list(string)
  default     = ["nexus", "kafka", "bookstore-ci"]
}

variable "internet_ports" {
  description = "TCP destination ports internet_namespaces may reach"
  type        = list(number)
  default     = [80, 443]
}

variable "enforce" {
  description = <<-EOT
    Create the final Deny policy. false installs only the Allow policies so
    their selectors can be checked against live traffic (Calico's
    default-allow applies to everything they miss); flip to true to enforce.
  EOT
  type        = bool
  default     = false
}
