SERVICE_NAME := catalog
VERSION := 1.0.0

CONTAINER_ENGINE ?= podman

APP_NAME := $(SERVICE_NAME)-app
COMMIT := $(shell git rev-parse --short HEAD)
PLATFORM := $(shell uname -s)
NAME_FLAG := -X github.com/codejsha/bookstore-microservices/$(SERVICE_NAME)/cmd.Name=$(APP_NAME)
VERSION_FLAG := -X github.com/codejsha/bookstore-microservices/$(SERVICE_NAME)/cmd.Version=$(VERSION)
COMMIT_FLAG := -X github.com/codejsha/bookstore-microservices/$(SERVICE_NAME)/cmd.Commit=$(COMMIT)

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
	@if [ -f ../../deploy/kubernetes/helm/$(SERVICE_NAME)/Chart.yaml ]; then \
		yq -i '.appVersion = "$(VERSION)"' ../../deploy/kubernetes/helm/$(SERVICE_NAME)/Chart.yaml; \
	fi

.PHONY: replace
replace:
	go mod edit -replace github.com/codejsha/shared-library-go=$(LIB_PATH)

.PHONY: deps
deps:
	go mod tidy
	go mod download

.PHONY: style
style:
	goimports -w ./internal/{application,domain,infrastructure}
	gofmt -w ./internal/{application,domain,infrastructure}

.PHONY: openapi-bundle
openapi-bundle:
	cd api/openapi && pnpm install --frozen-lockfile && pnpm exec tsp compile ./main.tsp --output-dir .

.PHONY: openapi-gen
openapi-gen:
	rm -rf ./generated/application/port/openapi
	oapi-codegen-cli -c api/openapi/config.toml

.PHONY: openapi
openapi: openapi-bundle openapi-gen

.PHONY: gorm
gorm:
	rm -rf ./generated/infrastructure/port/{dao,entity,repo}
	gorm-codegen-cli generate -e gen

.PHONY: api-gen
api-gen: spec openapi

.PHONY: gen
gen: api-gen gorm

.PHONY: build
build: spec
ifeq ($(PLATFORM),Darwin)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "$(NAME_FLAG) $(VERSION_FLAG) $(COMMIT_FLAG)" -o ./bin/$(APP_NAME)
else ifeq ($(PLATFORM),Linux)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(NAME_FLAG) $(VERSION_FLAG) $(COMMIT_FLAG)" -o ./bin/$(APP_NAME)
else
	@echo "Unsupported platform: $(PLATFORM)"
endif

.PHONY: build-linux
build-linux: spec
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(NAME_FLAG) $(VERSION_FLAG) $(COMMIT_FLAG)" -o ./bin/$(APP_NAME)

.PHONY: podman-build
podman-build: build-linux
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
	go test -v -json ./internal/...

.PHONY: run
run:
	go run ./main.go

.PHONY: migrate-up
migrate-up:
	$(CONTAINER_ENGINE) run --rm --network host -v $(PWD)/db/migrations:/migrations:ro \
		docker.io/migrate/migrate:v4.18.2 \
		-path=/migrations \
		-database "postgres://$${DB_USER:-postgres}:$${DB_PASSWORD}@$${DB_HOST:-localhost}:$${DB_PORT:-5432}/$(SERVICE_NAME)_db?sslmode=$${DB_SSLMODE:-disable}" \
		up

.PHONY: migrate-down
migrate-down:
	$(CONTAINER_ENGINE) run --rm --network host -v $(PWD)/db/migrations:/migrations:ro \
		docker.io/migrate/migrate:v4.18.2 \
		-path=/migrations \
		-database "postgres://$${DB_USER:-postgres}:$${DB_PASSWORD}@$${DB_HOST:-localhost}:$${DB_PORT:-5432}/$(SERVICE_NAME)_db?sslmode=$${DB_SSLMODE:-disable}" \
		down 1

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
