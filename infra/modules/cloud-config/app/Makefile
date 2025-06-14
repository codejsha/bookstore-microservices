SERVICE_NAME := config-server
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
	### helm chart
	yq -i '.appVersion = "$(VERSION)"' ../helm/$(SERVICE_NAME)/Chart.yaml
	### gradle
	perl -pi -e 's/version = ".*"/version=$(VERSION)/' gradle.properties
	### application
	yq -i '.spring.application.name = $(SERVICE_NAME)' src/main/resources/application.yaml
	yq -i '.spring.application.version = $(VERSION)' src/main/resources/application.yaml

.PHONY: clean
clean:
	sh gradlew clean

.PHONY: style
style:
	google-java-format -i $(shell find src -name "*.java")

.PHONY: gen
gen: spec style

.PHONY: build
build: spec
	sh gradlew build

.PHONY: docker-build
docker-build:
	docker buildx build -f ./ci/Dockerfile \
		--build-arg SERVICE_NAME=$(SERVICE_NAME) \
		-t harbor.example.com/bookstore/$(SERVICE_NAME):$(VERSION) .

.PHONY: docker-push
docker-push:
	docker image push harbor.example.com/bookstore/$(SERVICE_NAME):$(VERSION)

.PHONY: run
run:
	sh gradlew bootRun

.PHONY: test
test:
	sh gradlew test
