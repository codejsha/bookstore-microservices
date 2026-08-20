ARG STRIMZI_VERSION=0.51.0
ARG KAFKA_VERSION=4.2.0

FROM alpine:3.20 AS fetch

ARG DEBEZIUM_VERSION=3.0.7.Final
ARG AVRO_CONVERTER_VERSION=7.7.1

ARG DEBEZIUM_SHA512=""
ARG AVRO_SHA512=""

RUN apk add --no-cache curl unzip

WORKDIR /plugins

RUN set -eux; \
    url="https://repo1.maven.org/maven2/io/debezium/debezium-connector-postgres/${DEBEZIUM_VERSION}/debezium-connector-postgres-${DEBEZIUM_VERSION}-plugin.tar.gz"; \
    curl -fsSL -o /tmp/debezium.tgz "$url"; \
    if [ -n "$DEBEZIUM_SHA512" ]; then echo "$DEBEZIUM_SHA512  /tmp/debezium.tgz" | sha512sum -c -; fi; \
    tar -xzf /tmp/debezium.tgz -C /plugins; \
    rm /tmp/debezium.tgz

RUN set -eux; \
    url="https://hub-downloads.confluent.io/api/plugins/confluentinc/kafka-connect-avro-converter/versions/${AVRO_CONVERTER_VERSION}/confluentinc-kafka-connect-avro-converter-${AVRO_CONVERTER_VERSION}.zip"; \
    curl -fsSL -o /tmp/avro.zip "$url"; \
    if [ -n "$AVRO_SHA512" ]; then echo "$AVRO_SHA512  /tmp/avro.zip" | sha512sum -c -; fi; \
    unzip -q /tmp/avro.zip -d /plugins; \
    rm /tmp/avro.zip

FROM quay.io/strimzi/kafka:${STRIMZI_VERSION}-kafka-${KAFKA_VERSION}

USER root
COPY --from=fetch /plugins/ /opt/kafka/plugins/
RUN chown -R 1001:0 /opt/kafka/plugins && chmod -R g+rw /opt/kafka/plugins
USER 1001
