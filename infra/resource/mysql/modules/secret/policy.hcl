path "kv-bookstore/data/order/mysql*" { capabilities = ["read"] }
path "kv-bookstore/data/payment/mysql*" { capabilities = ["read"] }
path "kv-bookstore/data/delivery/mysql*" { capabilities = ["read"] }
path "kv-bookstore/data/notification/mysql*" { capabilities = ["read"] }
path "kv-bookstore/data/support/mysql*" { capabilities = ["read"] }
path "kv-bookstore/data/settlement/mysql*" { capabilities = ["read"] }

path "kv-infra/data/hyperswitch/payment" { capabilities = ["read"] }

path "database/creds/order-mysql-dynamic" { capabilities = ["read"] }
path "database/creds/payment-mysql-dynamic" { capabilities = ["read"] }
path "database/creds/delivery-mysql-dynamic" { capabilities = ["read"] }
path "database/creds/notification-mysql-dynamic" { capabilities = ["read"] }
path "database/creds/support-mysql-dynamic" { capabilities = ["read"] }
path "database/creds/settlement-mysql-dynamic" { capabilities = ["read"] }

path "database/static-creds/order-mysql-static" { capabilities = ["read"] }
path "database/static-creds/payment-mysql-static" { capabilities = ["read"] }
path "database/static-creds/delivery-mysql-static" { capabilities = ["read"] }
path "database/static-creds/notification-mysql-static" { capabilities = ["read"] }
path "database/static-creds/support-mysql-static" { capabilities = ["read"] }
path "database/static-creds/settlement-mysql-static" { capabilities = ["read"] }

path "database/static-creds/settlement-payment-mysql-static" { capabilities = ["read"] }
