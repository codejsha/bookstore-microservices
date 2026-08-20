SERVICE_NAME := notification
VERSION := 1.0.0

.PHONY: none
none:
	@echo "No target specified."

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: spec
spec:
	yq -i '.appVersion = "$(VERSION)"' ../../deploy/kubernetes/helm/$(SERVICE_NAME)/Chart.yaml

.PHONY: deps
deps:
	uv sync

.PHONY: style
style:
	uv run ruff check --fix .
	uv run ruff format .

.PHONY: openapi-bundle
openapi-bundle:
	cd api/openapi && pnpm install --frozen-lockfile && pnpm exec tsp compile ./main.tsp --output-dir .

.PHONY: openapi
openapi: openapi-bundle

.PHONY: gen
gen: spec openapi

.PHONY: build
build:
	uv build

.PHONY: podman-build
podman-build: build
	podman build \
		--no-cache \
		--platform linux/amd64 \
		--file ./ci/Dockerfile \
		--tag harbor.example.com/bookstore/$(SERVICE_NAME):$(VERSION) .

.PHONY: podman-push
podman-push: podman-build
	podman push harbor.example.com/bookstore/$(SERVICE_NAME):$(VERSION)

.PHONY: test
test:
	uv run pytest -v

.PHONY: run
run:
	NOTIFICATION__SERVER__MODE=debug uv run python main.py

.PHONY: migrate-create
migrate-create:
	uv run alembic revision --autogenerate -m "$(VAR)"

.PHONY: migrate-up
migrate-up:
	uv run alembic upgrade head

.PHONY: migrate-down
migrate-down:
	uv run alembic downgrade -1

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
