path "kv-bookstore/data/identity/*" { capabilities = ["read"] }
path "kv-bookstore/metadata/identity/*" { capabilities = ["read"] }
path "kv-infra/data/keycloak/admin/identity-admin/credentials" { capabilities = ["read"] }
path "kv-infra/metadata/keycloak/admin/identity-admin/credentials" { capabilities = ["read"] }
path "database/static-creds/identity-postgres-static" { capabilities = ["read"] }
