terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
    vault = {
      source  = "hashicorp/vault"
      version = "~> 5.10"
    }
    keycloak = {
      source  = "keycloak/keycloak"
      version = "~> 5.9"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.9"
    }
  }
}

ephemeral "vault_kv_secret_v2" "keycloak" {
  mount = "kv-infra"
  name  = "keycloak/admin/master-admin/credentials"
}

provider "keycloak" {
  client_id = "admin-cli"
  username  = ephemeral.vault_kv_secret_v2.keycloak.data["username"]
  password  = ephemeral.vault_kv_secret_v2.keycloak.data["password"]
  url       = var.keycloak_url
}

module "realm" {
  source     = "./modules/realm"
  realm_name = var.realm_name
  providers = {
    keycloak = keycloak
  }
}

data "keycloak_role" "master_admin" {
  realm_id = "master"
  name     = "admin"
}

module "master_admin" {
  source             = "./modules/master-admin"
  username           = var.master_admin_username
  email              = var.master_admin_email
  admin_role_id      = data.keycloak_role.master_admin.id
  keycloak_namespace = var.keycloak_namespace
  providers = {
    keycloak = keycloak
    vault    = vault
  }
  depends_on = [
    module.realm,
    module.realm_roles,
    module.argocd_oidc,
    module.identity_client,
    module.bookstore_admins,
    module.web_client,
    module.mobile_client,
    module.oauth2_proxy_client,
    module.audience_scope,
    module.temporal_worker_clients,
  ]
}

module "realm_roles" {
  source   = "./modules/realm-roles"
  realm_id = module.realm.realm_id
  roles = {
    USER    = "Signed-in customer: own profile, orders and public catalog"
    STAFF   = "Operations staff: catalog, inventory, delivery and support handling"
    MANAGER = "Full back-office access including users, payments and settlement"
    SYSTEM  = "Internal service-to-service automation"
  }
  providers = {
    keycloak = keycloak
  }
}

module "temporal_worker_clients" {
  source   = "./modules/temporal-worker-clients"
  realm_id = module.realm.realm_id
  services = ["order", "payment", "inventory", "notification", "delivery"]
  providers = {
    keycloak = keycloak
    vault    = vault
  }
}

module "infra_realm" {
  source             = "./modules/infra-realm"
  realm_name         = var.infra_realm_name
  bootstrap_accounts = var.infra_bootstrap_accounts
  roles = {
    DEVELOPER = "Developer: develop, commit and deploy applications and view monitoring; no infrastructure changes"
    MANAGER   = "Platform manager: may change infrastructure"
    ADMIN     = "Platform administrator: full control of infrastructure tools"
  }
  role_composites = {
    DEVELOPER = [module.temporal_oidc.permission_role_ids["bookstore:read"], module.temporal_oidc.permission_role_ids["temporal-system:read"]]
    MANAGER   = [module.temporal_oidc.permission_role_ids["bookstore:write"], module.temporal_oidc.permission_role_ids["temporal-system:read"]]
    ADMIN     = [module.temporal_oidc.permission_role_ids["temporal-system:admin"]]
  }
  providers = {
    keycloak = keycloak
    vault    = vault
  }
}

module "argocd_oidc" {
  source            = "./modules/argocd-oidc"
  realm_id          = module.infra_realm.realm_id
  groups_scope_name = module.infra_realm.groups_scope_name
  argocd_url        = var.argocd_url
  namespace         = var.argocd_namespace
  providers = {
    keycloak   = keycloak
    vault      = vault
    kubernetes = kubernetes
  }
}

module "kiali_oidc" {
  source            = "./modules/kiali-oidc"
  realm_id          = module.infra_realm.realm_id
  groups_scope_name = module.infra_realm.groups_scope_name
  kiali_url         = var.kiali_url
  namespace         = var.kiali_namespace
  providers = {
    keycloak   = keycloak
    vault      = vault
    kubernetes = kubernetes
  }
}

module "harbor_oidc" {
  source            = "./modules/harbor-oidc"
  realm_id          = module.infra_realm.realm_id
  groups_scope_name = module.infra_realm.groups_scope_name
  harbor_url        = var.harbor_url
  providers = {
    keycloak = keycloak
    vault    = vault
  }
}

moved {
  from = module.argocd_oidc.keycloak_openid_client_scope.groups
  to   = module.infra_realm.keycloak_openid_client_scope.groups
}

moved {
  from = module.argocd_oidc.keycloak_openid_user_realm_role_protocol_mapper.argocd_groups
  to   = module.infra_realm.keycloak_openid_user_realm_role_protocol_mapper.groups
}

module "temporal_oidc" {
  source            = "./modules/temporal-oidc"
  realm_id          = module.infra_realm.realm_id
  groups_scope_name = module.infra_realm.groups_scope_name
  temporal_url      = var.temporal_url
  namespace         = var.temporal_namespace
  providers = {
    keycloak   = keycloak
    vault      = vault
    kubernetes = kubernetes
  }
}

module "opensearch_oidc" {
  source            = "./modules/opensearch-oidc"
  realm_id          = module.infra_realm.realm_id
  groups_scope_name = module.infra_realm.groups_scope_name
  opensearch_url    = var.opensearch_url
  namespace         = var.opensearch_namespace
  providers = {
    keycloak   = keycloak
    vault      = vault
    kubernetes = kubernetes
  }
}

module "gitea_oidc" {
  source            = "./modules/gitea-oidc"
  realm_id          = module.infra_realm.realm_id
  groups_scope_name = module.infra_realm.groups_scope_name
  gitea_url         = var.gitea_url
  namespace         = var.gitea_namespace
  providers = {
    keycloak   = keycloak
    vault      = vault
    kubernetes = kubernetes
  }
}

module "grafana_oidc" {
  source      = "./modules/grafana-oidc"
  realm_id    = module.infra_realm.realm_id
  grafana_url = var.grafana_url
  namespace   = var.grafana_namespace
  providers = {
    keycloak   = keycloak
    vault      = vault
    kubernetes = kubernetes
  }
}

module "identity_client" {
  source   = "./modules/identity-client"
  realm_id = module.realm.realm_id
  providers = {
    keycloak = keycloak
    vault    = vault
  }
}

module "bookstore_admins" {
  source   = "./modules/bookstore-admins"
  realm_id = module.realm.realm_id
  accounts = var.bookstore_admin_accounts
  role_ids = module.realm_roles.role_ids
  providers = {
    keycloak = keycloak
  }
}

module "web_client" {
  source                          = "./modules/web-client"
  realm_id                        = module.realm.realm_id
  root_url                        = var.web_root_url
  valid_redirect_uris             = var.web_valid_redirect_uris
  valid_post_logout_redirect_uris = var.web_valid_post_logout_redirect_uris
  web_origins                     = var.web_web_origins
  providers = {
    keycloak = keycloak
  }
}

module "mobile_client" {
  source                          = "./modules/mobile-client"
  realm_id                        = module.realm.realm_id
  valid_redirect_uris             = var.mobile_valid_redirect_uris
  valid_post_logout_redirect_uris = var.mobile_valid_post_logout_redirect_uris
  providers = {
    keycloak = keycloak
  }
}

module "oauth2_proxy_client" {
  source                          = "./modules/oauth2-proxy-client"
  realm_id                        = module.realm.realm_id
  client_id                       = var.oauth2_proxy_client_id
  root_url                        = var.oauth2_proxy_root_url
  valid_redirect_uris             = var.oauth2_proxy_valid_redirect_uris
  valid_post_logout_redirect_uris = var.oauth2_proxy_valid_post_logout_redirect_uris
  web_origins                     = var.oauth2_proxy_web_origins
  providers = {
    keycloak = keycloak
    vault    = vault
  }
}

module "audience_scope" {
  source   = "./modules/audience-scope"
  realm_id = module.realm.realm_id
  audience = "bookstore"
  clients = {
    web          = module.web_client.client_uuid
    mobile       = module.mobile_client.client_uuid
    oauth2_proxy = module.oauth2_proxy_client.client_uuid
    identity     = module.identity_client.client_uuid
  }
  service_account_clients = ["identity"]
  providers = {
    keycloak = keycloak
  }
}

moved {
  from = module.identity_client.random_password.manager
  to   = module.bookstore_admins.random_password.account["devopsadmin"]
}

moved {
  from = module.identity_client.keycloak_user.manager
  to   = module.bookstore_admins.keycloak_user.account["devopsadmin"]
}

moved {
  from = module.identity_client.keycloak_user_roles.manager
  to   = module.bookstore_admins.keycloak_user_roles.account["devopsadmin"]
}

removed {
  from = module.identity_client.vault_kv_secret_v2.identity_keycloak

  lifecycle {
    destroy = false
  }
}
