SERVICE_NAME := commonlib-go
VERSION := 0.1.0

APP_NAME := $(SERVICE_NAME)
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
	### gradle
	perl -pi -e 's/version = ".*"/version=$(VERSION)/' gradle.properties

.PHONY: style
style:
	goimports -w ./internal/application/port/openapi
	gofmt -w ./internal/application/port/openapi

.PHONY: gen
gen: spec style

.PHONY: test
test:
	go test -v ./...
