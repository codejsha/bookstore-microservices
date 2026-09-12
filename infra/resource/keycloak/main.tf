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
  mount = "kv"
  name  = "keycloak/admin/credentials"
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
    module.web_client,
    module.mobile_client,
    module.oauth2_proxy_client,
    module.audience_scope,
  ]
}

module "realm_roles" {
  source   = "./modules/realm-roles"
  realm_id = module.realm.realm_id
  roles = {
    USER   = "Signed-in customer: own profile, orders and public catalog"
    STAFF  = "Operations staff: catalog, inventory, delivery and support handling"
    MANAGE = "Full back-office access including users, payments and settlement"
    SYSTEM = "Internal service-to-service automation"
  }
  providers = {
    keycloak = keycloak
  }
}

module "argocd_oidc" {
  source     = "./modules/argocd-oidc"
  realm_id   = module.realm.realm_id
  argocd_url = var.argocd_url
  namespace  = var.argocd_namespace
  providers = {
    keycloak   = keycloak
    vault      = vault
    kubernetes = kubernetes
  }
}

module "identity_client" {
  source               = "./modules/identity-client"
  realm_id             = module.realm.realm_id
  realm_admin_username = var.identity_realm_admin_username
  manage_role_id       = module.realm_roles.role_ids["MANAGE"]
  providers = {
    keycloak = keycloak
    vault    = vault
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
