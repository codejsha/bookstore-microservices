PROJECT_ROOT := $(shell cd ../../ && pwd)

.PHONY: none
none:
	@echo "No target specified."

.PHONY: grpc-bundle
grpc-bundle:
	tsp compile ./main.tsp --output-dir ./proto

.PHONY: grpc-gen
grpc-gen:
	rm -rf $(PROJECT_ROOT)/services/customer/generated/application/port/pb
	rm -rf $(PROJECT_ROOT)/services/identity/generated/application/port/pb
	rm -rf $(PROJECT_ROOT)/services/order/src/main/generated/com/codejsha/bookstore/generated/application/port/pb
	rm -rf $(PROJECT_ROOT)/services/payment/src/main/generated/com/codejsha/bookstore/generated/application/port/pb
	cd proto && buf lint
	cd proto && buf generate --template buf.gen.user.yaml
	cd proto && buf generate --template buf.gen.order.yaml
	cd proto && buf generate --template buf.gen.payment.yaml
	cd proto && buf generate --template buf.gen.delivery.yaml

.PHONY: grpc
grpc: grpc-bundle grpc-gen
