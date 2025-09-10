path "kv/data/bookstore/catalog/*" { capabilities = ["read"] }
path "kv/metadata/bookstore/catalog/*" { capabilities = ["read"] }
path "kv/data/opensearch/admin/credentials" { capabilities = ["read"] }
path "database/static-creds/catalog-postgres-static" { capabilities = ["read"] }
