terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
  }
}

locals {
  service_policies = {
    argocd = [
      { path = "kv-infra/data/argocd/admin/credentials", capabilities = ["create", "read", "update", "list"] },
      { path = "kv-infra/data/argocd/admin/credentials/*", capabilities = ["create", "read", "update", "list"] },
      { path = "kv-infra/data/gitea/ssh/host", capabilities = ["read"] },
    ]
    gitea = [
      { path = "kv-infra/data/gitea/*", capabilities = ["create", "read", "update", "list"] },
      { path = "pki_int/issuer/*", capabilities = ["read"] },
      { path = "pki_int/roles/gitea", capabilities = ["create", "read", "update", "delete"] },
      { path = "auth/kubernetes/role/gitea-issuer", capabilities = ["create", "read", "update", "delete"] },
      { path = "sys/policies/acl/gitea", capabilities = ["create", "read", "update", "delete"] },
      { path = "auth/kubernetes/role/gitea-role", capabilities = ["create", "read", "update", "delete"] },
    ]
    harbor = [
      { path = "kv-infra/data/harbor/*", capabilities = ["create", "read", "update", "list"] },
      { path = "pki_int/issuer/*", capabilities = ["read"] },
      { path = "pki_int/roles/harbor", capabilities = ["create", "read", "update", "delete"] },
      { path = "auth/kubernetes/role/harbor-issuer", capabilities = ["create", "read", "update", "delete"] },
      { path = "sys/policies/acl/vso-harbor", capabilities = ["create", "read", "update", "delete"] },
      { path = "auth/kubernetes/role/vso-harbor-role", capabilities = ["create", "read", "update", "delete"] },
    ]
    nexus = [
      { path = "kv-infra/data/nexus/*", capabilities = ["create", "read", "update", "list"] },
      { path = "pki_int/issuer/*", capabilities = ["read"] },
      { path = "pki_int/roles/nexus", capabilities = ["create", "read", "update", "delete"] },
      { path = "auth/kubernetes/role/nexus-issuer", capabilities = ["create", "read", "update", "delete"] },
      { path = "sys/policies/acl/nexus", capabilities = ["create", "read", "update", "delete"] },
      { path = "auth/kubernetes/role/nexus-role", capabilities = ["create", "read", "update", "delete"] },
    ]
    hyperswitch = [
      { path = "kv-infra/data/hyperswitch/*", capabilities = ["create", "read", "update", "list"] },
      { path = "sys/policies/acl/vso-hyperswitch", capabilities = ["create", "read", "update", "delete"] },
      { path = "auth/kubernetes/role/vso-hyperswitch-role", capabilities = ["create", "read", "update", "delete"] },
    ]
    keycloak = [
      { path = "kv-infra/data/keycloak/*", capabilities = ["create", "read", "update", "list"] },
      { path = "pki_int/issuer/*", capabilities = ["read"] },
      { path = "pki_int/roles/keycloak", capabilities = ["create", "read", "update", "delete"] },
      { path = "auth/kubernetes/role/keycloak-issuer", capabilities = ["create", "read", "update", "delete"] },
    ]
    oauth2-proxy = [
      { path = "sys/policies/acl/vso-oauth2-proxy", capabilities = ["create", "read", "update", "delete"] },
      { path = "auth/kubernetes/role/vso-oauth2-proxy-role", capabilities = ["create", "read", "update", "delete"] },
    ]
    opensearch = [
      { path = "kv-infra/data/opensearch/*", capabilities = ["create", "read", "update", "list"] },
      { path = "sys/policies/acl/opensearch", capabilities = ["create", "read", "update", "delete"] },
      { path = "auth/kubernetes/role/opensearch-role", capabilities = ["create", "read", "update", "delete"] },
      { path = "sys/policies/acl/vso-opensearch", capabilities = ["create", "read", "update", "delete"] },
      { path = "auth/kubernetes/role/vso-opensearch-role", capabilities = ["create", "read", "update", "delete"] },
    ]
    grafana = [
      { path = "kv-infra/data/seaweedfs/s3/credentials", capabilities = ["read"] },
      { path = "kv-infra/data/grafana/*", capabilities = ["create", "read", "update", "list"] },
      { path = "sys/policies/acl/vso-grafana", capabilities = ["create", "read", "update", "delete"] },
      { path = "auth/kubernetes/role/vso-grafana-role", capabilities = ["create", "read", "update", "delete"] },
    ]
    seaweedfs = [
      { path = "kv-infra/data/seaweedfs/*", capabilities = ["create", "read", "update", "list"] },
    ]
    prometheus = [
      { path = "kv-infra/data/opensearch/admin/credentials", capabilities = ["read"] },
    ]
    postgresql = [
      { path = "kv-infra/data/postgresql/*", capabilities = ["create", "read", "update", "list"] },
    ]
    kafka = [
      { path = "kv-infra/data/harbor/robots/pull/credentials", capabilities = ["read"] },
    ]
    temporal = [
      { path = "kv-infra/data/temporal/mysql", capabilities = ["read"] },
    ]
    flink             = []
    foundation-alerts = []
    argo-rollouts     = []
    kyverno           = []
  }

  config_policies = {
    argocd-config = [
      { path = "kv-infra/data/argocd/admin/credentials", capabilities = ["read"] },
      { path = "kv-infra/data/argocd/*", capabilities = ["create", "read", "update", "list"] },
      { path = "kv-infra/data/gitea/admin/credentials", capabilities = ["read"] },
      { path = "kv-infra/data/gitea/ssh/*", capabilities = ["read"] },
    ]
    gitea-config = [
      { path = "kv-infra/data/gitea/*", capabilities = ["create", "read", "update", "list"] },
      { path = "sys/policies/acl/gitea-ssh", capabilities = ["create", "read", "update", "delete"] },
      { path = "auth/kubernetes/role/gitea-ssh-role", capabilities = ["create", "read", "update", "delete"] },
      { path = "sys/policies/acl/config-server-gitea", capabilities = ["create", "read", "update", "delete"] },
      { path = "auth/kubernetes/role/config-server-gitea-role", capabilities = ["create", "read", "update", "delete"] },
    ]
    gitea-repo-config = [
      { path = "kv-infra/data/gitea/admin/credentials", capabilities = ["read"] },
    ]
    flink-config = [
      { path = "kv-infra/data/opensearch/admin/credentials", capabilities = ["read"] },
      { path = "kv-infra/data/harbor/robots/pull/credentials", capabilities = ["read"] },
      { path = "kv-infra/data/seaweedfs/s3/credentials", capabilities = ["read"] },
    ]
    hyperswitch-config = [
      { path = "kv-infra/data/hyperswitch/app/credentials", capabilities = ["read"] },
      { path = "kv-infra/data/hyperswitch/payment", capabilities = ["create", "read", "update"] },
    ]
    keycloak-config = [
      { path = "kv-infra/data/keycloak/*", capabilities = ["create", "read", "update", "list"] },
      { path = "kv-bookstore/data/admin/*", capabilities = ["create", "read", "update", "list"] },
      { path = "kv-bookstore/data/oauth2-proxy/*", capabilities = ["create", "read", "update", "list"] },
      { path = "kv-bookstore/data/order/temporal", capabilities = ["create", "read", "update", "list"] },
      { path = "kv-bookstore/data/payment/temporal", capabilities = ["create", "read", "update", "list"] },
      { path = "kv-bookstore/data/inventory/temporal", capabilities = ["create", "read", "update", "list"] },
      { path = "kv-bookstore/data/notification/temporal", capabilities = ["create", "read", "update", "list"] },
      { path = "kv-bookstore/data/delivery/temporal", capabilities = ["create", "read", "update", "list"] },
    ]
    kyverno-config = [
      { path = "kv-infra/data/harbor/robots/pull/credentials", capabilities = ["read"] },
      { path = "transit/keys/cosign-key", capabilities = ["read"] },
    ]
    harbor-config = [
      { path = "kv-infra/data/harbor/admin/credentials", capabilities = ["read"] },
      { path = "kv-infra/data/harbor/users/*", capabilities = ["create", "read", "update", "list"] },
      { path = "kv-infra/data/harbor/robots/*", capabilities = ["create", "read", "update", "list"] },
      { path = "kv-infra/data/harbor/ci/*", capabilities = ["create", "read", "update", "list"] },
      { path = "kv-infra/data/keycloak/harbor-oidc/client-secret", capabilities = ["read"] },
    ]
    nexus-config = [
      { path = "kv-infra/data/nexus/admin/credentials", capabilities = ["read"] },
      { path = "kv-infra/data/nexus/publisher", capabilities = ["create", "read", "update", "list"] },
    ]
    mysql-config = [
      { path = "kv-bookstore/data/+/mysql", capabilities = ["create", "read", "update", "list"] },
      { path = "sys/policies/acl/bookstore-mysql", capabilities = ["create", "read", "update", "delete"] },
      { path = "auth/kubernetes/role/bookstore-mysql-role", capabilities = ["create", "read", "update", "delete"] },
      { path = "database/config/*", capabilities = ["create", "read", "update", "delete"] },
      { path = "database/roles/*", capabilities = ["create", "read", "update", "delete"] },
      { path = "database/static-roles/*", capabilities = ["create", "read", "update", "delete"] },
      { path = "database/rotate-root/*", capabilities = ["create", "update"] },
      { path = "database/reset/*", capabilities = ["create", "update"] },
    ]
    opensearch-config = [
      { path = "kv-infra/data/opensearch/admin/credentials", capabilities = ["read"] },
    ]
    postgresql-config = [
      { path = "kv-bookstore/data/+/postgres", capabilities = ["create", "read", "update", "list"] },
      { path = "kv-bookstore/data/catalog/postgres-debezium", capabilities = ["create", "read", "update", "list"] },
      { path = "sys/policies/acl/bookstore-postgres", capabilities = ["create", "read", "update", "delete"] },
      { path = "auth/kubernetes/role/bookstore-postgres-role", capabilities = ["create", "read", "update", "delete"] },
      { path = "sys/policies/acl/bookstore-identity", capabilities = ["create", "read", "update", "delete"] },
      { path = "auth/kubernetes/role/bookstore-identity-role", capabilities = ["create", "read", "update", "delete"] },
      { path = "sys/policies/acl/bookstore-catalog", capabilities = ["create", "read", "update", "delete"] },
      { path = "auth/kubernetes/role/bookstore-catalog-role", capabilities = ["create", "read", "update", "delete"] },
      { path = "database/config/*", capabilities = ["create", "read", "update", "delete"] },
      { path = "database/roles/*", capabilities = ["create", "read", "update", "delete"] },
      { path = "database/static-roles/*", capabilities = ["create", "read", "update", "delete"] },
      { path = "database/rotate-root/*", capabilities = ["create", "update"] },
      { path = "database/reset/*", capabilities = ["create", "update"] },
    ]
    bookstore-base-config = [
      { path = "kv-infra/data/harbor/robots/pull/credentials", capabilities = ["read"] },
    ]
    tekton-config = [
      { path = "kv-infra/data/argocd/dev-ci/token", capabilities = ["read"] },
      { path = "kv-infra/data/gitea/admin/credentials", capabilities = ["read"] },
      { path = "kv-infra/data/gitea/resolver/token", capabilities = ["read"] },
      { path = "kv-infra/data/gitea/ssh/*", capabilities = ["read"] },
      { path = "kv-infra/data/gitea/webhook/*", capabilities = ["read"] },
      { path = "kv-infra/data/harbor/ci/credentials", capabilities = ["read"] },
      { path = "kv-infra/data/nexus/publisher", capabilities = ["read"] },
      { path = "transit/keys/cosign-key", capabilities = ["create", "read", "update", "delete"] },
      { path = "transit/keys/cosign-key/config", capabilities = ["create", "read", "update"] },
      { path = "transit/sign/cosign-key/*", capabilities = ["create", "read", "update"] },
      { path = "transit/verify/cosign-key/*", capabilities = ["create", "read", "update"] },
      { path = "sys/policies/acl/tekton-chains-sign", capabilities = ["create", "read", "update", "delete"] },
      { path = "auth/kubernetes/role/tekton-chains-role", capabilities = ["create", "read", "update", "delete"] },
    ]
  }

  password_generators = toset([
    "argocd", "gitea", "grafana", "harbor", "hyperswitch", "keycloak", "nexus", "opensearch",
    "gitea-config", "harbor-config", "keycloak-config", "mysql-config", "postgresql-config",
  ])
  password_generate_rules = [
    { path = "sys/policies/password/password-special/generate", capabilities = ["read"] },
    { path = "sys/policies/password/password-alphanumeric/generate", capabilities = ["read"] },
  ]

  grafana_alert_readers = toset([
    "harbor", "argocd", "opensearch", "gitea", "kafka",
    "prometheus", "keycloak", "flink", "nexus",
    "foundation-alerts", "argo-rollouts",
    "kyverno",
  ])
  grafana_alert_read_rules = [
    { path = "kv-infra/data/grafana/admin/credentials", capabilities = ["read"] },
  ]

  all_policies = {
    for name, rules in merge(local.service_policies, local.config_policies) :
    name => concat(
      rules,
      contains(local.password_generators, name) ? local.password_generate_rules : [],
      contains(local.grafana_alert_readers, name) ? local.grafana_alert_read_rules : [],
    )
  }
}

locals {
  kv_data_pattern = "^kv-(bookstore|infra)/data/"

  policies_expanded = {
    for name, rules in local.all_policies : name => flatten([
      for r in rules : (
        can(regex(local.kv_data_pattern, r.path))
        ? flatten([
          for base in [r.path] : (
            length(setintersection(toset(r.capabilities), toset(["create", "update", "delete"]))) > 0
            ? [
              { path = base, capabilities = distinct(concat(r.capabilities, ["delete"])) },
              { path = replace(base, "/\\/data\\//", "/metadata/"), capabilities = ["create", "read", "update", "list", "delete"] },
            ]
            : [
              { path = base, capabilities = r.capabilities },
              { path = replace(base, "/\\/data\\//", "/metadata/"), capabilities = ["read"] },
            ]
          )
        ])
        : [r]
      )
    ])
  }
}

resource "vault_policy" "service" {
  for_each = local.policies_expanded
  name     = "${each.key}-policy"
  policy = join("\n", [
    for rule in each.value : <<-EOT
path "${rule.path}" {
  capabilities = ${jsonencode(rule.capabilities)}
}
EOT
  ])
}

resource "vault_auth_backend" "approle" {
  type = "approle"
}

resource "vault_approle_auth_backend_role" "service" {
  for_each       = local.all_policies
  backend        = vault_auth_backend.approle.path
  role_name      = each.key
  token_policies = [vault_policy.service[each.key].name]
  token_ttl      = 3600
  token_max_ttl  = 7200
}

data "vault_approle_auth_backend_role_id" "service" {
  for_each  = local.all_policies
  backend   = vault_auth_backend.approle.path
  role_name = vault_approle_auth_backend_role.service[each.key].role_name
}

resource "vault_approle_auth_backend_role_secret_id" "service" {
  for_each  = local.all_policies
  backend   = vault_auth_backend.approle.path
  role_name = vault_approle_auth_backend_role.service[each.key].role_name
}
