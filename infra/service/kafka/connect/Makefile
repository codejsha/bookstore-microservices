CONTAINER_ENGINE ?= podman

IMAGE                  ?= harbor.example.com/bookstore/kafka-connect:4.2.0-dbz3.0.7
STRIMZI_VERSION        ?= 0.51.0
KAFKA_VERSION          ?= 4.2.0
DEBEZIUM_VERSION       ?= 3.0.7.Final
AVRO_CONVERTER_VERSION ?= 7.7.1

DEBEZIUM_SHA512 ?=
AVRO_SHA512     ?=

.PHONY: none
none:
	@echo "No target specified."

.PHONY: help
help:
	@echo "Targets:"
	@echo "  image   Build the Kafka Connect worker image"
	@echo "  push    Push the image to the registry"
	@echo ""
	@echo "Vars: IMAGE, STRIMZI_VERSION, KAFKA_VERSION, DEBEZIUM_VERSION,"
	@echo "      AVRO_CONVERTER_VERSION, DEBEZIUM_SHA512, AVRO_SHA512"

.PHONY: image
image:
	$(CONTAINER_ENGINE) build -t $(IMAGE) \
	  --build-arg STRIMZI_VERSION=$(STRIMZI_VERSION) \
	  --build-arg KAFKA_VERSION=$(KAFKA_VERSION) \
	  --build-arg DEBEZIUM_VERSION=$(DEBEZIUM_VERSION) \
	  --build-arg AVRO_CONVERTER_VERSION=$(AVRO_CONVERTER_VERSION) \
	  --build-arg DEBEZIUM_SHA512=$(DEBEZIUM_SHA512) \
	  --build-arg AVRO_SHA512=$(AVRO_SHA512) \
	  -f docker/Dockerfile .

.PHONY: push
push:
	$(CONTAINER_ENGINE) push $(IMAGE)
