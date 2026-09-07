SERVICE_NAME := settlement
VERSION := 1.0.0
GRADLE_VERSION ?= 9.4.1

CONTAINER_ENGINE ?= podman
GRADLE ?= sh gradlew

APP_NAME := $(SERVICE_NAME)-app
COMMIT := $(shell git rev-parse --short HEAD)
PLATFORM := $(shell uname -s)

.PHONY: none
none:
	@echo "No target specified."

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: commit
commit:
	@echo $(COMMIT)

.PHONY: spec
spec:
	yq -i '.appVersion = "$(VERSION)"' ../../deploy/kubernetes/helm/$(SERVICE_NAME)/Chart.yaml
	perl -pi -e 's/^version\s*=.*/version=$(VERSION)/' gradle.properties
	yq -i '.spring.application.name = "$(SERVICE_NAME)"' src/main/resources/application.yaml
	yq -i '.spring.application.version = "$(VERSION)"' src/main/resources/application.yaml

.PHONY: wrapper
wrapper:
	perl -pi -e 's,gradle-[0-9.]+-bin\.zip,gradle-$(GRADLE_VERSION)-bin.zip,' gradle/wrapper/gradle-wrapper.properties
	gradle wrapper --gradle-version $(GRADLE_VERSION) --console=plain
	$(GRADLE) wrapper --gradle-version $(GRADLE_VERSION) --console=plain

.PHONY: clean
clean:
	$(GRADLE) clean --console=plain

.PHONY: style
style:
	$(GRADLE) spotlessApply --console=plain

.PHONY: openapi-bundle
openapi-bundle:
	cd api/openapi && pnpm install --frozen-lockfile && pnpm exec tsp compile ./main.tsp --output-dir .

.PHONY: openapi-gen
openapi-gen:
	rm -rf ./src/main/generated/com/codejsha/bookstore/generated/application/port/openapi
	oapi-codegen-cli -c api/openapi/config.toml

.PHONY: openapi
openapi: openapi-bundle openapi-gen

.PHONY: jooq
jooq:
	rm -rf ./src/main/generated/com/codejsha/bookstore/generated/infrastructure/adapter/jooq
	$(GRADLE) jooqGenerateContainer --console=plain

.PHONY: api-gen
api-gen: spec openapi

.PHONY: gen
gen: api-gen jooq

.PHONY: build
build: spec
	$(GRADLE) build --info --console=plain

.PHONY: podman-build
podman-build: build
	podman build \
		--no-cache \
		--platform linux/amd64 \
		--file ./ci/Dockerfile \
		--build-arg SERVICE_NAME=$(SERVICE_NAME) \
		--tag harbor.example.com/bookstore/$(SERVICE_NAME):$(VERSION) .

.PHONY: podman-push
podman-push: podman-build
	podman push harbor.example.com/bookstore/$(SERVICE_NAME):$(VERSION)

.PHONY: test
test:
	$(GRADLE) test --console=plain

.PHONY: run
run:
	$(GRADLE) bootRun --console=plain

.PHONY: migrate-create
migrate-create:
	@if [ -z "$(VAR)" ]; then echo "Usage: make migrate-create VAR=name"; exit 1; fi
	@next=$$(ls db/migrations/V*__*.sql 2>/dev/null | sed -E 's|.*/V([0-9]+)__.*|\1|' | sort -n | tail -1); \
	next=$$((next + 1)); \
	file="db/migrations/V$${next}__$(VAR).sql"; \
	: > $$file; \
	echo "Created $$file"

.PHONY: migrate-up
migrate-up:
	$(CONTAINER_ENGINE) run --rm --network host -v $(PWD)/db/migrations:/flyway/sql:ro \
		docker.io/flyway/flyway:10-alpine \
		-url=jdbc:mysql://$${DB_HOST:-localhost}:$${DB_PORT:-3306}/$(SERVICE_NAME)_db \
		-user=$${DB_USER:-root} -password=$${DB_PASSWORD} \
		migrate

.PHONY: migrate-info
migrate-info:
	$(CONTAINER_ENGINE) run --rm --network host -v $(PWD)/db/migrations:/flyway/sql:ro \
		docker.io/flyway/flyway:10-alpine \
		-url=jdbc:mysql://$${DB_HOST:-localhost}:$${DB_PORT:-3306}/$(SERVICE_NAME)_db \
		-user=$${DB_USER:-root} -password=$${DB_PASSWORD} \
		info

CONTRACT_SCHEMA_PATH := $(abspath api/openapi/schema/openapi.1.0.0.json)

.PHONY: contract-test-deps
contract-test-deps:
	cd ci/contract-test && uv sync

.PHONY: contract-test
contract-test:
	cd ci/contract-test && CONTRACT_SCHEMA_PATH=$(CONTRACT_SCHEMA_PATH) uv run pytest

.PHONY: contract-test-cli
contract-test-cli:
	cd ci/contract-test && uv run schemathesis run \
		--url $${CONTRACT_BASE_URL:-http://localhost:8080} \
		--checks all \
		$(CONTRACT_SCHEMA_PATH)

.PHONY: load-test-smoke
load-test-smoke:
	cd ci/load-test && k6 run smoke.js

.PHONY: load-test
load-test:
	cd ci/load-test && k6 run load.js
