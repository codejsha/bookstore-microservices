path "kv-bookstore/data/identity/*" { capabilities = ["read"] }
path "kv-bookstore/metadata/identity/*" { capabilities = ["read"] }
path "kv-bookstore/data/keycloak/admin/credentials" { capabilities = ["read"] }
path "kv-bookstore/metadata/keycloak/admin/credentials" { capabilities = ["read"] }
path "database/static-creds/identity-postgres-static" { capabilities = ["read"] }
