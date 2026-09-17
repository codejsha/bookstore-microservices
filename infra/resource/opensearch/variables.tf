variable "opensearch_api_url" {
  description = "OpenSearch REST API URL (reachable from where Terraform runs)"
  type        = string
}

variable "keycloak_issuer_url" {
  description = "Keycloak platform-infra realm issuer URL as published in its OIDC discovery document"
  type        = string
}

variable "oidc_client_id" {
  description = "Keycloak client ID of OpenSearch Dashboards, required as the ID token audience"
  type        = string
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
