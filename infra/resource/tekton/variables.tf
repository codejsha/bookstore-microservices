variable "namespace" {
  description = "Namespace name"
  type        = string
}

variable "service_accounts" {
  description = "Service account base names (suffixed with -sa). Each carries the registry pull Secret, which is where Tekton Chains gets the credentials it pushes signatures with."
  type        = list(string)
}

variable "harbor_address" {
  description = "Harbor address"
  type        = string
}

variable "gitea_address" {
  description = "Gitea address"
  type        = string
}

variable "gitea_org" {
  description = "Gitea organization that owns the source repos (for webhook provisioning)"
  type        = string
}

variable "source_repo_suffix" {
  description = "Suffix on the Gitea source repo name. Must match the gitea-repo unit's source_repo_suffix — a mismatch points the webhooks at repos that do not exist."
  type        = string
}

variable "resolver_namespace" {
  description = "Namespace of the Tekton remote-resolvers deployment. The operator installs it into targetNamespace, not the upstream release's separate tekton-pipelines-resolvers namespace"
  type        = string
  default     = "tekton-pipelines"
}

variable "chains_namespace" {
  description = "Namespace the operator installs the Chains controller into"
  type        = string
  default     = "tekton-pipelines"
}

variable "chains_service_account" {
  description = "ServiceAccount the Chains controller runs as; bound to the Vault tekton-chains-role"
  type        = string
  default     = "tekton-chains-controller"
}

variable "registry_pull_secret" {
  description = "dockerconfigjson Secret attached to the CI ServiceAccounts so Chains can push signatures and attestations to Harbor"
  type        = string
  default     = "harbor-ci-pull"
}

variable "resolver_token_secret" {
  description = "Secret in the resolver namespace holding the Gitea API token. The v1.9.5 resolver image ships no ssh binary, so remote resources are fetched over the Gitea SCM API instead of an SSH clone"
  type        = string
  default     = "gitea-resolver-token"
}

variable "vault_url" {
  description = "Vault URL"
  type        = string
}

variable "vault_auth_role" {
  type = string
}

variable "vault_k8s_jwt" {
  description = "Short-lived Kubernetes ServiceAccount token for Vault kubernetes-auth login"
  type        = string
  sensitive   = true
}

variable "uv_services" {
  description = "List of uv (Python) service names for Tekton pipelines"
  type        = list(string)
}

variable "golang_services" {
  description = "List of Go service names for Tekton pipelines"
  type        = list(string)
}

variable "gradle_services" {
  description = "List of Kotlin service names for Tekton pipelines"
  type        = list(string)
}

variable "vite_services" {
  description = "List of Vite SPA names for Tekton pipelines"
  type        = list(string)
}

variable "harbor_registry" {
  description = <<-EOT
    Registry endpoint CI pushes to. This must be the same hostname pods reference: Tekton Chains
    stamps it into the in-toto subject, and an attestation whose subject names the in-cluster
    endpoint will not verify against an image admitted under the external one.
  EOT
  type        = string
  default     = "harbor.example.com"
}

variable "vault_internal_url" {
  description = "In-cluster Vault address. vault_url is the external one the terraform provider talks to; build pods only resolve cluster DNS."
  type        = string
  default     = "http://vault.vault.svc.cluster.local:8200"
}

variable "nexus_maven_url" {
  description = "Nexus maven group URL the Gradle builds resolve against"
  type        = string
}

variable "argocd_internal_address" {
  description = "In-cluster ArgoCD API address used by the CI tasks (build pods only resolve cluster DNS)"
  type        = string
  default     = "argocd-server.argocd.svc.cluster.local"
}

variable "nexus_pypi_url" {
  description = "Nexus pypi group simple-index URL the uv builds resolve against"
  type        = string
}

variable "nexus_raw_url" {
  description = "Nexus raw-hosted repo base (no trailing slash) codegen-verify-external / codegen-verify-grpc fetch the mirrored specs (idpapi/payapi, bookstore-grpc proto bundle) from"
  type        = string
}

variable "codegen_verify_targets" {
  description = "Codegen drift checks the golang/gradle/uv CI pipelines run: any of openapi (needs the oapi-codegen-cli image in Harbor), external (idpapi/payapi from the Nexus raw repo; identity and payment only), grpc (proto bundle from the Nexus raw repo; needs the protoc-toolchain image in Harbor; customer/identity/order/payment/delivery), db (needs the gorm-codegen-cli image and the privileged dind sidecar admitted in the CI namespace). Empty runs none."
  type        = list(string)
  default     = []
}

variable "nexus_npm_url" {
  description = "Nexus npm group URL the vite builds resolve against"
  type        = string
}

variable "nexus_docker_registry" {
  description = <<-EOT
    host:port of the Nexus docker connector (the ghcr.io pull-through cache).
    Unlike the other Nexus endpoints this is not a /repository/ URL — Nexus OSS
    serves each docker repo on its own port, so clients address the connector
    Service directly.
  EOT
  type        = string
}

variable "max_concurrent_pipelineruns" {
  description = "How many PipelineRuns may run at once. Enforced as a pod quota (2 pods per in-flight linear run); further runs queue as ExceededResourceQuota rather than failing."
  type        = number
  default     = 2
}

variable "eventlistener_count" {
  description = "Long-lived EventListener pods in the namespace — one per language template. They permanently hold a slot in the pod quota above."
  type        = number
  default     = 4
}
