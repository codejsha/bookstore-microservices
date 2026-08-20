PACKAGE_NAME := @bookstore/design
VERSION := 1.0.0
PUBLISH_SCRIPT := ../scripts/nexus-npm-publish.sh
VAULT_PUBLISHER_PATH := kv/nexus/publisher
VAULT_ADDR ?= http://vault.example.com
NEXUS_CA_FILE ?=

.PHONY: none
none:
	@echo "No target specified."

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: deps
deps:
	pnpm install

.PHONY: build
build:
	pnpm build

.PHONY: build-tokens
build-tokens:
	pnpm build:tokens

.PHONY: context
context:
	pnpm context

.PHONY: style
style:
	pnpm check

.PHONY: storybook
storybook:
	pnpm storybook

.PHONY: publish
publish: build
	@VAULT_ADDR="$(VAULT_ADDR)" \
	VAULT_PUBLISHER_PATH="$(VAULT_PUBLISHER_PATH)" \
	NEXUS_CA_FILE="$(NEXUS_CA_FILE)" \
	$(PUBLISH_SCRIPT) .
