SERVICE_NAME := config-server
VERSION := 1.0.0
GRADLE_VERSION ?= 9.4.1

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
	yq -i '.appVersion = "$(VERSION)"' ../helm/Chart.yaml
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

.PHONY: gen
gen: spec

.PHONY: build
build: spec
	sh gradlew build --info --console=plain

.PHONY: podman-build
podman-build: build
	podman build \
		--platform linux/amd64 \
		--file ./ci/Dockerfile \
		--build-arg SERVICE_NAME=$(SERVICE_NAME) \
		--tag harbor.example.com/bookstore/$(SERVICE_NAME):$(VERSION) .

.PHONY: podman-push
podman-push:
	podman push harbor.example.com/bookstore/$(SERVICE_NAME):$(VERSION)

.PHONY: skopeo-push
skopeo-push:
	skopeo copy docker-daemon:harbor.example.com/bookstore/$(SERVICE_NAME):$(VERSION) \
		docker://harbor.example.com/bookstore/$(SERVICE_NAME):$(VERSION) \
		--dest-cert-dir . --debug

.PHONY: test
test:
	sh gradlew test --console=plain

.PHONY: run
run:
	sh gradlew bootRun --console=plain
