SERVICE_NAME := order
VERSION := 0.1.0

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
	### openapi
	yq -i '.apis.core@v1.decorators.info-override.version = "$(VERSION)"' ./api/openapi/redocly.yaml
	### helm chart
	yq -i '.appVersion = "$(VERSION)"' ../../deploy/kubernetes/helm/$(SERVICE_NAME)-command/Chart.yaml
	yq -i '.appVersion = "$(VERSION)"' ../../deploy/kubernetes/helm/$(SERVICE_NAME)-query/Chart.yaml
	### gradle
	perl -pi -e 's/version = ".*"/version=$(VERSION)/' gradle.properties
	### application
	yq -i '.spring.application.name = "$(SERVICE_NAME)"' src/main/resources/application.yaml
	yq -i '.spring.application.version = "$(VERSION)"' src/main/resources/application.yaml

.PHONY: clean
clean:
	sh gradlew clean

.PHONY: style
style:
	sh gradlew spotlessApply

.PHONY: buf
buf:
	cd ../../api/proto && sh generate.sh

.PHONY: openapi-bundle
openapi-bundle:
	redocly bundle --config ./api/openapi/redocly.yaml core@v1

.PHONY: openapi-gen
openapi-gen:
	rm -rf ./src/main/generated/com/example/bookstore/service/application/port/openapi
	openapi-generator-cli generate -c ../../api/openapi/kotlin/config-all-reactive.yaml
	openapi-generator-cli generate -c ../../api/openapi/kotlin/config-params.yaml

.PHONY: openapi
openapi: openapi-bundle openapi-gen

.PHONY: gen
gen: spec buf openapi style

.PHONY: build
build: spec
	sh gradlew build --info

.PHONY: podman-build
podman-build:
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
	sh gradlew test

.PHONY: run
run:
	sh gradlew bootRun

.PHONY: migrate
migrate:
	sh gradlew flywayInfo
	sh gradlew flywayMigrate
