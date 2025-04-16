SERVICE_NAME := identity
VERSION := 0.1.0
KEYCLOAK_VERSION := 26.1.2

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
	### openapi
	yq -i '.apis.core@v1.decorators.info-override.version = "$(VERSION)"' ./api/openapi/redocly.yaml
	yq -i '.apis.oidc@v1.decorators.info-override.version = "$(VERSION)"' ./api/idp/redocly.yaml
	### helm chart
	yq -i '.appVersion = "$(VERSION)"' ../../deploy/kubernetes/helm/$(SERVICE_NAME)/Chart.yaml

.PHONY: replace
replace:
	go mod edit -replace=github.com/codejsha/bookstore-microservices/commonlib-go=$(LIB_REPO)@$(LIB_VERSION)

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: style
style:
	goimports -w ./internal/application/port/openapi
	gofmt -w ./internal/application/port/openapi
	goimports -w ./internal/application/port/idp
	gofmt -w ./internal/application/port/idp

.PHONY: buf
buf:
	rm -rf ./internal/application/port/pb
	cd ../../api/proto && bash generate.sh

.PHONY: openapi-bundle
openapi-bundle:
	redocly bundle --config ./api/openapi/redocly.yaml core@v1

.PHONY: openapi-gen
openapi-gen:
	rm -rf ./internal/application/port/openapi
	openapi-generator-cli generate -c ../../api/openapi/golang/config-all.yaml
	openapi-generator-cli generate -c ../../api/openapi/golang/config-params.yaml

.PHONY: openapi
openapi: openapi-bundle openapi-gen

.PHONY: idpapi-download
idpapi-download:
	curl -o ./api/idp/admin-openapi.json -L https://www.keycloak.org/docs-api/$(KEYCLOAK_VERSION)/rest-api/openapi.json

.PHONY: idpapi-bundle
idpapi-bundle:
	redocly bundle --config ./api/idp/redocly.yaml oidc@v1
	redocly bundle --config ./api/idp/redocly.yaml admin@v1

.PHONY: idpapi-gen
idpapi-gen:
	rm -rf ./internal/application/port/idp
	openapi-generator-cli generate -c ../../api/idp/config-oidc-all.yaml
	openapi-generator-cli generate -c ../../api/idp/config-oidc-params.yaml
	openapi-generator-cli generate -c ../../api/idp/config-admin-all.yaml
	openapi-generator-cli generate -c ../../api/idp/config-admin-params.yaml
	rm -rf ./internal/application/port/idp/.openapi-generator

.PHONY: idpapi
idpapi: idpapi-download idpapi-bundle idpapi-gen

.PHONY: gen
gen: version buf openapi idpapi style

.PHONY: wire
wire:
	wire check ./internal/di
	wire show ./internal/di
	wire gen ./internal/di

.PHONY: build
build: wire
ifeq ($(PLATFORM),Darwin)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "$(NAME_FLAG) $(VERSION_FLAG) $(COMMIT_FLAG)" -o ./bin/$(APP_NAME)
else ifeq ($(PLATFORM),Linux)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(NAME_FLAG) $(VERSION_FLAG) $(COMMIT_FLAG)" -o ./bin/$(APP_NAME)
else
	@echo "Unsupported platform: $(PLATFORM)"
endif

.PHONY: docker-build
docker-build:
	docker buildx build -f ./ci/package/Dockerfile \
		--platform linux/amd64,darwin/arm64 \
		--build-arg SERVICE_NAME=$(SERVICE_NAME) \
		-t harbor.harbor.svc.cluster.local/library/bookstore/$(SERVICE_NAME):$(VERSION) .

.PHONY: run
run:
	go run ./main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: migrate-create
migrate-create:
	### make migrate-create VAR=create_table
	migrate create -ext sql -dir db/migrations -seq $(VAR)
