variable "vault_url" {
  description = "Vault address for the kubernetes-auth provider login"
  type        = string
}

variable "vault_auth_role" {
  description = "Vault kubernetes-auth role"
  type        = string
}

variable "vault_k8s_jwt" {
  description = "Short-lived Kubernetes ServiceAccount token for Vault kubernetes-auth login"
  type        = string
  sensitive   = true
}

variable "org_name" {
  description = "Gitea organization that owns every pushed repo"
  type        = string
}

variable "gitea_api_url" {
  description = "Gitea API base URL for repo creation"
  type        = string
}

variable "gitea_http_url" {
  description = "Gitea HTTP base URL for git push"
  type        = string
}

variable "cloud_config_repo" {
  description = "Gitea repo name for the Spring Cloud Config content"
  type        = string
  default     = "cloud-config"
}

variable "cloud_config_dir" {
  description = "Absolute path to deploy/kubernetes/cloud-config/repo"
  type        = string
}

variable "app_repos" {
  description = "List of <svc>-helm repos to push charts for (must match the argocd unit's app_repos)"
  type        = list(string)
}

variable "helm_dir" {
  description = "Absolute path to deploy/kubernetes/helm"
  type        = string
}

variable "services" {
  description = "Service names under services/<name> whose source trees get mirrored to Gitea"
  type        = list(string)
}

variable "repo_root" {
  description = "Absolute repo root (get_repo_root); source is pushed via `git archive` against this working tree"
  type        = string
}

variable "source_repo_suffix" {
  description = "Suffix appended to the service name for its source repo."
  type        = string
  default     = "-source"
}

variable "tekton_catalog_repo" {
  description = "Gitea repo name for the mirrored Tekton catalog"
  type        = string
  default     = "tekton-catalog"
}

variable "tekton_catalog_upstream" {
  description = "Upstream Tekton catalog git URL to mirror"
  type        = string
  default     = "https://github.com/tektoncd/catalog"
}

variable "tekton_custom_catalog_repo" {
  description = "Gitea repo name for the in-repo custom Tekton catalog (pipelines + custom tasks)"
  type        = string
  default     = "tektoncd-custom-catalog"
}

variable "tekton_custom_catalog_dir" {
  description = "Absolute path to infra/resource/tekton/catalog"
  type        = string
}

variable "lib_repos" {
  description = "GitHub repos mirrored into Gitea under the same name; CI builds the services against their source rather than a published artifact"
  type        = set(string)
}

variable "github_owner" {
  description = "GitHub owner/org holding the upstream library repos"
  type        = string
}
