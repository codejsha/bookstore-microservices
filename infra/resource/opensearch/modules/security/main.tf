terraform {
  required_providers {
    opensearch = {
      source = "opensearch-project/opensearch"
    }
    restapi = {
      source = "Mastercard/restapi"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
  }
}

data "kubernetes_config_map_v1" "vault_pki_ca" {
  metadata {
    name      = var.ca_configmap_name
    namespace = var.ca_configmap_namespace
  }
}

locals {
  securityconfig_path = "/_plugins/_security/api/securityconfig"

  securityconfig = {
    filtered_alias_mode             = "warn"
    disable_rest_auth               = false
    disable_intertransport_auth     = false
    respect_request_indices_options = false
    kibana = {
      multitenancy_enabled   = true
      private_tenant_enabled = true
      default_tenant         = ""
      server_username        = "kibanaserver"
      index                  = ".kibana"
      sign_in_options        = ["BASIC", "OPENID"]
    }
    http = {
      anonymous_auth_enabled = false
      xff = {
        enabled         = false
        internalProxies = "192\\.168\\.0\\.10|192\\.168\\.0\\.11"
        remoteIpHeader  = "X-Forwarded-For"
      }
    }
    authc = {
      jwt_auth_domain = {
        http_enabled = false
        order        = 0
        http_authenticator = {
          challenge = false
          type      = "jwt"
          config = {
            jwks_uri                         = "https://your-jwks-endpoint.com/.well-known/jwks.json"
            signing_key                      = "base64 encoded HMAC key or public RSA/ECDSA pem key"
            jwt_header                       = "Authorization"
            jwt_clock_skew_tolerance_seconds = 30
          }
        }
        authentication_backend = {
          type   = "noop"
          config = {}
        }
        description = "Authenticate via Json Web Token"
      }
      openid_auth_domain = {
        http_enabled = true
        order        = 1
        http_authenticator = {
          challenge = false
          type      = "openid"
          config = {
            subject_key        = "preferred_username"
            roles_key          = "groups"
            openid_connect_url = "${var.keycloak_issuer_url}/.well-known/openid-configuration"
            required_audience  = var.oidc_client_id
            required_issuer    = var.keycloak_issuer_url
            openid_connect_idp = {
              enable_ssl            = true
              verify_hostnames      = true
              pemtrustedcas_content = data.kubernetes_config_map_v1.vault_pki_ca.data[var.ca_configmap_key]
            }
          }
        }
        authentication_backend = {
          type   = "noop"
          config = {}
        }
        description = "Authenticate Keycloak platform-infra ID tokens forwarded by OpenSearch Dashboards"
      }
      ldap = {
        http_enabled = false
        order        = 5
        http_authenticator = {
          challenge = false
          type      = "basic"
          config    = {}
        }
        authentication_backend = {
          type = "ldap"
          config = {
            enable_ssl             = false
            enable_start_tls       = false
            enable_ssl_client_auth = false
            verify_hostnames       = true
            hosts                  = ["localhost:8389"]
            userbase               = "ou=people,dc=example,dc=com"
            usersearch             = "(sAMAccountName={0})"
          }
        }
        description = "Authenticate via LDAP or Active Directory"
      }
      basic_internal_auth_domain = {
        http_enabled = true
        order        = 4
        http_authenticator = {
          challenge = true
          type      = "basic"
          config    = {}
        }
        authentication_backend = {
          type   = "intern"
          config = {}
        }
        description = "Authenticate via HTTP Basic against internal users database"
      }
      proxy_auth_domain = {
        http_enabled = false
        order        = 3
        http_authenticator = {
          challenge = false
          type      = "proxy"
          config = {
            user_header  = "x-proxy-user"
            roles_header = "x-proxy-roles"
          }
        }
        authentication_backend = {
          type   = "noop"
          config = {}
        }
        description = "Authenticate via proxy"
      }
      clientcert_auth_domain = {
        http_enabled = false
        order        = 2
        http_authenticator = {
          challenge = false
          type      = "clientcert"
          config = {
            username_attribute = "cn"
          }
        }
        authentication_backend = {
          type   = "noop"
          config = {}
        }
        description = "Authenticate via SSL client certificates"
      }
      kerberos_auth_domain = {
        http_enabled = false
        order        = 6
        http_authenticator = {
          challenge = true
          type      = "kerberos"
          config = {
            krb_debug                  = false
            strip_realm_from_principal = true
          }
        }
        authentication_backend = {
          type   = "noop"
          config = {}
        }
      }
    }
    authz = {
      roles_from_another_ldap = {
        http_enabled = false
        authorization_backend = {
          type   = "ldap"
          config = {}
        }
        description = "Authorize via another Active Directory"
      }
      roles_from_myldap = {
        http_enabled = false
        authorization_backend = {
          type = "ldap"
          config = {
            enable_ssl             = false
            enable_start_tls       = false
            enable_ssl_client_auth = false
            verify_hostnames       = true
            hosts                  = ["localhost:8389"]
            rolebase               = "ou=groups,dc=example,dc=com"
            rolesearch             = "(member={0})"
            userrolename           = "disabled"
            rolename               = "cn"
            resolve_nested_roles   = true
            userbase               = "ou=people,dc=example,dc=com"
            usersearch             = "(uid={0})"
          }
        }
        description = "Authorize via LDAP or Active Directory"
      }
    }
    auth_failure_listeners         = {}
    do_not_fail_on_forbidden       = false
    multi_rolespan_enabled         = true
    hosts_resolver_mode            = "ip-only"
    do_not_fail_on_forbidden_empty = false
    on_behalf_of = {
      enabled = false
    }
  }

  role_mappings = {
    all_access = {
      description   = "Maps admin to all_access"
      backend_roles = ["admin", "ADMIN"]
    }
    security_manager = {
      description   = null
      backend_roles = ["ADMIN"]
    }
    kibana_user = {
      description   = "Maps kibanauser to kibana_user"
      backend_roles = ["kibanauser", "DEVELOPER", "MANAGER", "ADMIN"]
    }
    readall = {
      description   = null
      backend_roles = ["readall", "DEVELOPER"]
    }
    platform_manager = {
      description   = null
      backend_roles = ["MANAGER"]
    }
  }
}

resource "restapi_object" "securityconfig" {
  path                    = local.securityconfig_path
  object_id               = "config"
  create_path             = "${local.securityconfig_path}/config"
  create_method           = "PUT"
  read_path               = local.securityconfig_path
  update_path             = "${local.securityconfig_path}/config"
  update_method           = "PUT"
  destroy_path            = local.securityconfig_path
  destroy_method          = "GET"
  ignore_server_additions = true
  data = jsonencode({
    "dynamic" = local.securityconfig
  })
}

resource "opensearch_role" "platform_manager" {
  role_name   = "platform_manager"
  description = "Platform manager: index data and management, snapshots, ISM and templates without security administration"

  cluster_permissions = [
    "cluster_monitor",
    "cluster_composite_ops",
    "manage_snapshots",
    "cluster:admin/opensearch/snapshot_management/*",
    "cluster:admin/opendistro/ism/*",
    "cluster_manage_index_templates",
  ]

  index_permissions {
    index_patterns  = ["*"]
    allowed_actions = ["crud", "create_index", "manage", "indices_monitor"]
  }

  tenant_permissions {
    tenant_patterns = ["global_tenant"]
    allowed_actions = ["kibana_all_write"]
  }
}

resource "opensearch_roles_mapping" "platform" {
  for_each      = local.role_mappings
  role_name     = each.key
  description   = each.value.description
  backend_roles = each.value.backend_roles
  depends_on    = [opensearch_role.platform_manager]
}
