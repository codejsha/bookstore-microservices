path "kv-bookstore/data/catalog/postgres*" { capabilities = ["read"] }
path "kv-bookstore/data/customer/postgres*" { capabilities = ["read"] }
path "kv-bookstore/data/identity/postgres*" { capabilities = ["read"] }
path "kv-bookstore/data/inventory/postgres*" { capabilities = ["read"] }

path "database/creds/*" { capabilities = ["read"] }

path "database/static-creds/*" { capabilities = ["read"] }
