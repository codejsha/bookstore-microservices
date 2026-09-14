terraform {
  required_providers {
    harbor = {
      source = "goharbor/harbor"
    }
  }
}

resource "harbor_config_auth" "oidc" {
  auth_mode         = "oidc_auth"
  primary_auth_mode = true

  oidc_name                     = "Keycloak"
  oidc_endpoint                 = var.oidc_endpoint
  oidc_client_id                = var.oidc_client_id
  oidc_client_secret_wo         = var.oidc_client_secret
  oidc_client_secret_wo_version = 1
  oidc_scope                    = "openid,profile,email,groups"
  oidc_verify_cert              = true
  oidc_auto_onboard             = true
  oidc_user_claim               = "preferred_username"
  oidc_groups_claim             = "groups"
  oidc_admin_group              = var.oidc_admin_group
}
