SERVICE_NAME := catalog-indexer
VERSION := 1.2.0

COMMIT := $(shell git rev-parse --short HEAD)

.PHONY: none
none:
	@echo "No target specified."

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: commit
commit:
	@echo $(COMMIT)

.PHONY: podman-build
podman-build:
	podman build \
		--no-cache \
		--platform linux/amd64 \
		--file ./Dockerfile \
		--tag harbor.example.com/bookstore/$(SERVICE_NAME):$(VERSION) .

.PHONY: podman-push
podman-push: podman-build
	podman push harbor.example.com/bookstore/$(SERVICE_NAME):$(VERSION)
