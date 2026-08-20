SERVICES := admin catalog customer delivery identity inventory order payment
CONFIG := config.toml
SPEC_DIR := ../../services
PUBLISH_SCRIPT := ../../scripts/nexus-npm-publish.sh
VAULT_PUBLISHER_PATH := kv/nexus/publisher
VAULT_ADDR ?= http://vault.example.com
NEXUS_CA_FILE ?=

.PHONY: none
none:
	@echo "No target specified."

.PHONY: generate
generate:
	@for svc in $(SERVICES); do \
		echo "==> generating @bookstore/$$svc-client"; \
		rm -rf $$svc/src; \
		oapi-codegen-cli -c $(CONFIG) \
			-p input_spec=$(SPEC_DIR)/$$svc/api/openapi/schema/openapi.1.0.0.json \
			-p output_dir=$$svc/src; \
	done

.PHONY: build
build: generate
	@for svc in $(SERVICES); do \
		echo "==> building @bookstore/$$svc-client"; \
		( cd $$svc && pnpm install && pnpm build ); \
	done

.PHONY: publish
publish:
	@VAULT_ADDR="$(VAULT_ADDR)" \
	VAULT_PUBLISHER_PATH="$(VAULT_PUBLISHER_PATH)" \
	NEXUS_CA_FILE="$(NEXUS_CA_FILE)" \
	$(PUBLISH_SCRIPT) $(SERVICES)
