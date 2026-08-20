include "root" {
  path = find_in_parent_folders("terragrunt.hcl")
}

include "kubernetes" {
  path = find_in_parent_folders("_envcommon/provider_kubernetes.hcl")
}

include "vault" {
  path = find_in_parent_folders("_envcommon/provider_vault.hcl")
}

dependencies {
  paths = [
    "../../service/keycloak",
    "../../service/vault",
  ]
}

locals {
  common     = read_terragrunt_config(find_in_parent_folders("_envcommon/vault_k8s.hcl")).locals
  vault_role = "keycloak-config"
}

inputs = {
  vault_url       = local.common.vault_url
  vault_auth_role = "tf-${local.vault_role}"
  vault_k8s_jwt = run_cmd("--terragrunt-quiet", "kubectl", "create", "token",
  "tf-${local.vault_role}", "-n", local.common.vault_ns, "--duration=20m")

  keycloak_url       = "https://keycloak.example.com"
  keycloak_namespace = "keycloak"

  master_admin_username = "devopsadmin"
  master_admin_email    = "devopsadmin@example.com"

  realm_name       = "bookstore"
  argocd_url       = "https://argocd.example.com"
  argocd_namespace = "argocd"

  web_root_url                        = "http://localhost:5173"
  web_valid_redirect_uris             = ["http://localhost:5173/*", "https://bookstore.example.com/*"]
  web_valid_post_logout_redirect_uris = ["http://localhost:5173/*", "https://bookstore.example.com/*"]
  web_web_origins                     = ["http://localhost:5173", "https://bookstore.example.com"]

  identity_realm_admin_username = "devopsadmin@example.com"

  mobile_valid_redirect_uris = [
    "bookstore://callback",
    "exp://*",
  ]
  mobile_valid_post_logout_redirect_uris = [
    "bookstore://logout",
    "exp://*",
  ]

  oauth2_proxy_client_id = "bookstore-edge"
  oauth2_proxy_root_url  = "https://bookstore.example.com"
  oauth2_proxy_valid_redirect_uris = [
    "https://bookstore.example.com/oauth2/callback",
    "https://bookstore-admin.example.com/oauth2/callback",
    "https://cloud-config.example.com/oauth2/callback",
  ]
  oauth2_proxy_valid_post_logout_redirect_uris = [
    "https://bookstore.example.com/",
    "https://bookstore-admin.example.com/",
    "https://cloud-config.example.com/",
  ]
  oauth2_proxy_web_origins = [
    "https://bookstore.example.com",
    "https://bookstore-admin.example.com",
    "https://cloud-config.example.com",
  ]
}

terraform {
  source = "${get_repo_root()}/infra//resource/keycloak"
}
