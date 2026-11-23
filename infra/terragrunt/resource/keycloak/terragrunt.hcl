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
    "../../resource/cert-manager-issuers",
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

  master_admin_username = "master-admin"
  master_admin_email    = "master-admin@example.com"

  realm_name = "bookstore"

  infra_realm_name = "platform-infra"
  infra_bootstrap_accounts = {
    administrator       = { username = "platform-admin", email = "platform-admin@example.com", first_name = "Platform", last_name = "Admin", roles = ["PLATFORM_ADMIN"], store_in_vault = true }
    operator            = { username = "platform-devops", email = "platform-devops@example.com", first_name = "Platform", last_name = "DevOps", roles = ["PLATFORM_DEVOPS"] }
    developer           = { username = "platform-dev", email = "platform-dev@example.com", first_name = "Platform", last_name = "Dev", roles = ["PLATFORM_DEV"] }
    bookstore_admin     = { username = "bookstore-admin", email = "bookstore-admin@example.com", first_name = "Bookstore", last_name = "Admin", roles = ["BOOKSTORE_ADMIN"] }
    bookstore_operator  = { username = "bookstore-devops", email = "bookstore-devops@example.com", first_name = "Bookstore", last_name = "DevOps", roles = ["BOOKSTORE_DEVOPS"] }
    bookstore_developer = { username = "bookstore-dev", email = "bookstore-dev@example.com", first_name = "Bookstore", last_name = "Dev", roles = ["BOOKSTORE_DEV"] }
  }
  argocd_url           = "https://argocd.example.com"
  argocd_namespace     = "argocd"
  kiali_url            = "https://kiali.example.com"
  kiali_namespace      = "istio-system"
  grafana_url          = "https://grafana.example.com"
  harbor_url           = "https://harbor.example.com"
  gitea_url            = "https://git.example.com"
  gitea_namespace      = "gitea"
  temporal_url         = "https://temporal.example.com"
  temporal_namespace   = "temporal"
  opensearch_url       = "https://opensearch.example.com"
  opensearch_namespace = "opensearch"
  grafana_namespace    = "grafana"

  web_root_url                        = "http://localhost:5173"
  web_valid_redirect_uris             = ["http://localhost:5173/*", "https://bookstore.example.com/*"]
  web_valid_post_logout_redirect_uris = ["http://localhost:5173/*", "https://bookstore.example.com/*"]
  web_web_origins                     = ["http://localhost:5173", "https://bookstore.example.com"]

  bookstore_admin_accounts = {
    manager = { username = "bookstore-manager", email = "bookstore-manager@example.com", first_name = "Bookstore", last_name = "Manager", role = "MANAGER" }
  }

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
  ]
  oauth2_proxy_valid_post_logout_redirect_uris = [
    "https://bookstore.example.com/",
    "https://bookstore-admin.example.com/",
  ]
  oauth2_proxy_web_origins = [
    "https://bookstore.example.com",
    "https://bookstore-admin.example.com",
  ]
}

terraform {
  source = "${get_repo_root()}/infra//resource/keycloak"
}
