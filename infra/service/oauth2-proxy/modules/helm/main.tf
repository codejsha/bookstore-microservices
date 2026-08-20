terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}

locals {
  base_args = {
    "provider"            = "oidc"
    "oidc-issuer-url"     = var.oidc_issuer_url
    "skip-oidc-discovery" = "true"
    "login-url"           = var.login_url
    "redeem-url"          = var.redeem_url
    "oidc-jwks-url"       = var.oidc_jwks_url
    "redirect-url"        = var.redirect_url
    "scope"               = var.oidc_scope

    "pass-access-token"         = "true"
    "set-authorization-header"  = "true"
    "pass-authorization-header" = "true"
    "set-xauthrequest"          = "true"

    "skip-jwt-bearer-tokens" = "true"

    "reverse-proxy"         = "true"
    "cookie-refresh"        = var.cookie_refresh
    "cookie-secure"         = "true"
    "cookie-samesite"       = "lax"
    "cookie-domain"         = var.cookie_domains
    "whitelist-domain"      = var.whitelist_domain
    "code-challenge-method" = "S256"
    "skip-provider-button"  = "true"
  }

  session_values = var.session_store == "redis" ? {
    sessionStorage = {
      type = "redis"
      redis = {
        clientType = "standalone"
        standalone = {
          connectionUrl = var.redis_url
        }
      }
    }
  } : {}

  values = merge({
    config = {
      existingSecret = var.existing_secret
      emailDomains   = ["*"]
      upstreams      = ["static://200"]
      cookieName     = "_bookstore_edge"
    }
    service = {
      portNumber  = var.service_port
      appProtocol = "http"
    }
    extraArgs = local.base_args
  }, local.session_values)
}

resource "helm_release" "oauth2_proxy" {
  namespace  = var.namespace
  name       = "oauth2-proxy"
  repository = "https://oauth2-proxy.github.io/manifests"
  chart      = "oauth2-proxy"
  version    = var.chart_version

  values = [yamlencode(local.values)]
}
