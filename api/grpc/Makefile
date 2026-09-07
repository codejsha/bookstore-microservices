PROJECT_ROOT := $(shell cd ../../ && pwd)
GRPC_API_VERSION := 1.0.0
SERVICES := customer identity order payment delivery

TOOLCHAIN_IMAGE := protoc-toolchain
TOOLCHAIN_VERSION := 0.1.0
REGISTRY ?= harbor.example.com/bookstore
PROTOC_VERSION := 35.1
BUF_VERSION := 1.72.0
GRPC_JAVA_VERSION := 1.78.0
PROTOC_GEN_GO_VERSION := v1.36.11
PROTOC_GEN_GO_GRPC_VERSION := v1.6.2

NEXUS_RAW_URL ?= https://nexus.example.com/repository/raw-hosted
NEXUS_AUTH ?=
nexus_auth = { if [ -n "$(NEXUS_AUTH)" ]; then printf '%s' "$(NEXUS_AUTH)"; \
	else u=$$(vault kv get -mount=kv -field=username nexus/publisher) && \
	     p=$$(vault kv get -mount=kv -field=password nexus/publisher) && \
	     [ -n "$$u" ] && [ -n "$$p" ] && printf '%s:%s' "$$u" "$$p"; fi; }

GRPC_PROTO_URL := $(NEXUS_RAW_URL)/vendor/bookstore-grpc/$(GRPC_API_VERSION)/proto.tar.gz

.PHONY: none
none:
	@echo "No target specified."

.PHONY: grpc-bundle
grpc-bundle:
	pnpm install --frozen-lockfile && pnpm exec tsp compile ./main.tsp --output-dir ./proto
	cd proto && buf lint

.PHONY: grpc-gen
grpc-gen:
	rm -rf $(PROJECT_ROOT)/services/customer/generated/application/port/pb
	rm -rf $(PROJECT_ROOT)/services/identity/generated/application/port/pb
	rm -rf $(PROJECT_ROOT)/services/order/src/main/generated/com/codejsha/bookstore/generated/application/port/pb
	rm -rf $(PROJECT_ROOT)/services/payment/src/main/generated/com/codejsha/bookstore/generated/application/port/pb
	cd proto && buf generate --template buf.gen.user.yaml
	cd proto && buf generate --template buf.gen.order.yaml
	cd proto && buf generate --template buf.gen.payment.yaml
	cd proto && buf generate --template buf.gen.delivery.yaml
	@for svc in $(SERVICES); do \
		dest=../../services/$$svc/api/grpc/proto; \
		rm -rf $$dest && mkdir -p $$dest && cp -R proto/buf.yaml proto/modules $$dest/; \
	done
	$(MAKE) --no-print-directory -C ../../services/delivery grpc-gen

.PHONY: grpc
grpc: grpc-bundle grpc-gen

.PHONY: grpc-mirror
grpc-mirror: grpc-bundle
	@auth=$$($(nexus_auth)); \
	[ -n "$$auth" ] || { echo "ERROR: no Nexus auth (set NEXUS_AUTH or ensure Vault kv/nexus/publisher is reachable)" >&2; exit 1; }; \
	bundle=$$(mktemp) && tar -czf "$$bundle" -C proto buf.yaml modules && \
	curl -fsS -u "$$auth" --upload-file "$$bundle" "$(GRPC_PROTO_URL)"; rc=$$?; rm -f "$$bundle"; exit $$rc
	@echo "mirrored -> $(GRPC_PROTO_URL)"

.PHONY: podman-build
podman-build:
	podman build \
		--platform linux/amd64 \
		--build-arg PROTOC_VERSION=$(PROTOC_VERSION) \
		--build-arg BUF_VERSION=$(BUF_VERSION) \
		--build-arg GRPC_JAVA_VERSION=$(GRPC_JAVA_VERSION) \
		--build-arg PROTOC_GEN_GO_VERSION=$(PROTOC_GEN_GO_VERSION) \
		--build-arg PROTOC_GEN_GO_GRPC_VERSION=$(PROTOC_GEN_GO_GRPC_VERSION) \
		--file ./toolchain/Dockerfile \
		--tag $(REGISTRY)/$(TOOLCHAIN_IMAGE):$(TOOLCHAIN_VERSION) ./toolchain

.PHONY: podman-versions
podman-versions:
	podman run --rm --platform linux/amd64 $(REGISTRY)/$(TOOLCHAIN_IMAGE):$(TOOLCHAIN_VERSION) sh -c \
		'buf --version; protoc --version; protoc-gen-go --version; protoc-gen-go-grpc --version'

.PHONY: podman-push
podman-push: podman-build
	podman push $(REGISTRY)/$(TOOLCHAIN_IMAGE):$(TOOLCHAIN_VERSION)
