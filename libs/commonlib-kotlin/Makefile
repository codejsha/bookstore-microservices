SERVICE_NAME := commonlib-kotlin
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

.PHONY: clean
clean:
	sh gradlew clean

.PHONY: style
style:
	sh gradlew spotlessApply

.PHONY: gen
gen: spec style

.PHONY: build
build: spec
	sh gradlew build --info

.PHONY: test
test:
	sh gradlew test

.PHONY: publish
publishing:
	sh gradlew publish
