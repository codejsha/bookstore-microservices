SERVICE_NAME := admin
VERSION := 1.0.0
GRADLE_VERSION ?= 9.4.1

APP_NAME := $(SERVICE_NAME)-app
COMMIT := $(shell git rev-parse --short HEAD)
PLATFORM := $(shell uname -s)

DOWNSTREAMS := catalog customer identity inventory order payment settlement

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
	sh gradlew wrapper --gradle-version $(GRADLE_VERSION) --console=plain

.PHONY: clean
clean:
	sh gradlew clean --console=plain

.PHONY: style
style:
	sh gradlew spotlessApply --console=plain

.PHONY: openapi-bundle
openapi-bundle:
	cd api/openapi && pnpm install --frozen-lockfile && pnpm exec tsp compile ./main.tsp --output-dir .

.PHONY: openapi-gen
openapi-gen:
	rm -rf ./src/main/generated/com/codejsha/bookstore/generated/application/port/openapi
	oapi-codegen-cli -c api/openapi/config.toml

.PHONY: openapi
openapi: openapi-bundle openapi-gen

.PHONY: clients
clients:
	rm -rf ./src/main/generated/com/codejsha/bookstore/generated/application/port/restclient
	@for svc in $(DOWNSTREAMS); do \
		echo "generating $$svc client"; \
		oapi-codegen-cli -c api/clients/$$svc/config.toml || exit 1; \
	done

.PHONY: gen
gen: spec openapi clients

.PHONY: build
build: spec
	sh gradlew build --info --console=plain

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
	sh gradlew test --console=plain

.PHONY: run
run:
	sh gradlew bootRun --console=plain

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
