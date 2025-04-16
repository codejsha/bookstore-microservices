#!/bin/bash
trap 'echo "${bash_source[0]}: line ${lineno}: status ${?}: user ${user}: func ${funcname[0]}"' err
set -o errexit
set -o errtrace
set -o xtrace

buf ls-files

rm -rf ../../services/catalog/internal/application/port/pb
rm -rf ../../services/customer/internal/application/port/pb
rm -rf ../../services/identity/internal/application/port/pb
rm -rf ../../services/inventory/internal/application/port/pb
rm -rf ../../services/order/src/main/generated/com/example/bookstore/service/application/port/pb
rm -rf ../../services/payment/src/main/generated/com/example/bookstore/service/application/port/pb

buf generate --template buf.gen.author.yaml
buf generate --template buf.gen.book.yaml
buf generate --template buf.gen.category.yaml
buf generate --template buf.gen.order.yaml
buf generate --template buf.gen.payment.yaml
buf generate --template buf.gen.publisher.yaml
buf generate --template buf.gen.stock.yaml
buf generate --template buf.gen.user.yaml
