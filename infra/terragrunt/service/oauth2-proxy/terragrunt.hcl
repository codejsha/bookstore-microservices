include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "kubernetes" {
  path = find_in_parent_folders("_envcommon/provider_kubernetes.hcl")
}

include "helm" {
  path = find_in_parent_folders("_envcommon/provider_helm.hcl")
}

include "vault" {
  path = find_in_parent_folders("_envcommon/provider_vault.hcl")
}

dependencies {
  paths = [
    "../../service/keycloak",
    "../../service/istio",
    "../../service/gateway-api",
    "../../service/vault",
    "../../service/vault-secrets-operator",
    "../../resource/bookstore-base",
    "../../resource/keycloak",
    "../../resource/valkey",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role = "oauth2-proxy"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  namespace = "bookstore"

  oidc_issuer_url = "https://keycloak.example.com/realms/bookstore"
  login_url       = "https://keycloak.example.com/realms/bookstore/protocol/openid-connect/auth"
  redeem_url      = "http://keycloak-service.keycloak.svc.cluster.local:8080/realms/bookstore/protocol/openid-connect/token"
  oidc_jwks_url   = "http://keycloak-service.keycloak.svc.cluster.local:8080/realms/bookstore/protocol/openid-connect/certs"
  redirect_url    = "/oauth2/callback"
  oidc_scope      = "openid email profile roles"

  cookie_domains   = ".example.com"
  whitelist_domain = ".example.com"

  session_store = "redis"
  redis_url     = "redis://oauth2-proxy-valkey-primary.bookstore.svc.cluster.local:6379"

  app_host = "bookstore.example.com"

  admin_app_host         = "bookstore-admin.example.com"
  admin_web_service_name = "admin-web"

  protected_hosts = ["cloud-config.example.com"]
}

terraform {
  source = "${get_repo_root()}/infra//service/oauth2-proxy"
}
