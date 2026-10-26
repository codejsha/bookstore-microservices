path "kv-bookstore/data/catalog/postgres*" { capabilities = ["read"] }
path "kv-bookstore/data/customer/postgres*" { capabilities = ["read"] }
path "kv-bookstore/data/identity/postgres*" { capabilities = ["read"] }
path "kv-bookstore/data/order/postgres*" { capabilities = ["read"] }
path "kv-bookstore/data/payment/postgres*" { capabilities = ["read"] }
path "kv-bookstore/data/fulfillment/postgres*" { capabilities = ["read"] }
path "kv-bookstore/data/settlement/postgres*" { capabilities = ["read"] }

path "kv-bookstore/data/order/temporal" { capabilities = ["read"] }
path "kv-bookstore/data/payment/temporal" { capabilities = ["read"] }
path "kv-bookstore/data/fulfillment/temporal" { capabilities = ["read"] }

path "kv-infra/data/hyperswitch/payment" { capabilities = ["read"] }

path "database/creds/*" { capabilities = ["read"] }

path "database/static-creds/*" { capabilities = ["read"] }
