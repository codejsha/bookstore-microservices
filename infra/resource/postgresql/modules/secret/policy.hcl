path "kv/data/bookstore/catalog/postgres*" { capabilities = ["read"] }
path "kv/data/bookstore/customer/postgres*" { capabilities = ["read"] }
path "kv/data/bookstore/identity/postgres*" { capabilities = ["read"] }
path "kv/data/bookstore/inventory/postgres*" { capabilities = ["read"] }

path "database/creds/*" { capabilities = ["read"] }

path "database/static-creds/*" { capabilities = ["read"] }
