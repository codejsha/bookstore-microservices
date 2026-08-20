SERVICE_NAME := delivery
VERSION := 1.0.0
PROTO_DIR := ../../api/grpc/proto/modules
PB_OUT := generated/application/port/pb/deliverypb

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

.PHONY: grpc
grpc:
	rm -rf $(PB_OUT)
	mkdir -p $(PB_OUT)
	uv run python -m grpc_tools.protoc \
		-I=$(PROTO_DIR) \
		--python_out=$(PB_OUT) \
		--grpc_python_out=$(PB_OUT) \
		--pyi_out=$(PB_OUT) \
		$(PROTO_DIR)/delivery/v1.proto
	@touch $(PB_OUT)/__init__.py $(PB_OUT)/delivery/__init__.py
	@sed -i.bak 's|^from delivery import v1_pb2|from . import v1_pb2|' \
		$(PB_OUT)/delivery/v1_pb2_grpc.py
	@rm -f $(PB_OUT)/delivery/v1_pb2_grpc.py.bak

.PHONY: openapi-bundle
openapi-bundle:
	cd api/openapi && pnpm install --frozen-lockfile && pnpm exec tsp compile ./main.tsp --output-dir .

.PHONY: openapi-gen
openapi-gen:
	rm -rf ./generated/application/port/openapi ./generated/application/port/model ./generated/application/port/api ./generated/application/port/constant
	oapi-codegen-cli -c api/openapi/config.toml
	@find ./generated/application/port/model -name "*.py" -exec sed -i.bak \
		-e 's|^from \.shipment_status import|from ..constant.shipment_status import|' \
		-e 's|^from \.carrier_status import|from ..constant.carrier_status import|' \
		-e 's|^from \.freight_status import|from ..constant.freight_status import|' \
		{} +
	@find ./generated/application/port/model -name "*.bak" -delete

.PHONY: openapi
openapi: openapi-bundle openapi-gen

.PHONY: gen
gen: spec grpc openapi

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
	DELIVERY__SERVER__MODE=debug uv run python main.py

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
