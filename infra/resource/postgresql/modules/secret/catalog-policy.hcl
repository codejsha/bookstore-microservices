path "kv-bookstore/data/catalog/*" { capabilities = ["read"] }
path "kv-bookstore/metadata/catalog/*" { capabilities = ["read"] }
path "kv-infra/data/opensearch/admin/credentials" { capabilities = ["read"] }
path "database/static-creds/catalog-postgres-static" { capabilities = ["read"] }
