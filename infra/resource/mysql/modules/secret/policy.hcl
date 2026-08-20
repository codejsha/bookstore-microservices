path "kv/data/bookstore/order/mysql*" { capabilities = ["read"] }
path "kv/data/bookstore/payment/mysql*" { capabilities = ["read"] }
path "kv/data/bookstore/delivery/mysql*" { capabilities = ["read"] }
path "kv/data/bookstore/notification/mysql*" { capabilities = ["read"] }
path "kv/data/bookstore/support/mysql*" { capabilities = ["read"] }
path "kv/data/bookstore/settlement/mysql*" { capabilities = ["read"] }
path "kv/data/bookstore/identity/keycloak" { capabilities = ["read"] }

path "kv/data/hyperswitch/payment" { capabilities = ["read"] }

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
