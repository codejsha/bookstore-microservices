ARG FLINK_VERSION=1.20.5

FROM docker.io/library/maven:3.9-eclipse-temurin-17 AS runner-build
ARG FLINK_VERSION
ARG OPENSEARCH_CONNECTOR_VERSION=2.0.0-1.19
WORKDIR /build
COPY runner/pom.xml ./
COPY runner/src ./src
RUN mvn -q -DskipTests -Dflink.version=${FLINK_VERSION} package

ADD https://repo1.maven.org/maven2/org/apache/flink/flink-sql-connector-opensearch2/${OPENSEARCH_CONNECTOR_VERSION}/flink-sql-connector-opensearch2-${OPENSEARCH_CONNECTOR_VERSION}.jar /build/flink-sql-connector-opensearch2-${OPENSEARCH_CONNECTOR_VERSION}.jar
RUN mkdir -p /build/patch/META-INF/services && \
    echo "org.apache.flink.opensearch2.shaded.org.opensearch.common.xcontent.spi.XContentProvider" \
      > "/build/patch/META-INF/services/org.apache.flink.opensearch2.shaded.org.opensearch.core.xcontent.spi.MediaTypeProvider" && \
    cd /build/patch && \
    jar -uf /build/flink-sql-connector-opensearch2-${OPENSEARCH_CONNECTOR_VERSION}.jar \
      "META-INF/services/org.apache.flink.opensearch2.shaded.org.opensearch.core.xcontent.spi.MediaTypeProvider"

FROM flink:${FLINK_VERSION}-java17

ARG FLINK_VERSION
ARG KAFKA_CONNECTOR_VERSION=3.4.0-1.20
ARG REDIS_CONNECTOR_VERSION=1.4.3
ARG OPENSEARCH_CONNECTOR_VERSION=2.0.0-1.19
ARG AVRO_CONFLUENT_VERSION=${FLINK_VERSION}

USER root

COPY --from=runner-build /build/target/flink-sql-runner.jar /opt/flink/lib/flink-sql-runner.jar
COPY --from=runner-build /build/flink-sql-connector-opensearch2-${OPENSEARCH_CONNECTOR_VERSION}.jar /opt/flink/lib/

ADD https://repo1.maven.org/maven2/org/apache/flink/flink-sql-connector-kafka/${KAFKA_CONNECTOR_VERSION}/flink-sql-connector-kafka-${KAFKA_CONNECTOR_VERSION}.jar /opt/flink/lib/
ADD https://repo1.maven.org/maven2/org/apache/flink/flink-sql-avro-confluent-registry/${AVRO_CONFLUENT_VERSION}/flink-sql-avro-confluent-registry-${AVRO_CONFLUENT_VERSION}.jar /opt/flink/lib/
ADD https://repo1.maven.org/maven2/io/github/jeff-zou/flink-connector-redis/${REDIS_CONNECTOR_VERSION}/flink-connector-redis-${REDIS_CONNECTOR_VERSION}.jar /opt/flink/lib/

RUN chmod 644 /opt/flink/lib/*.jar

RUN mkdir -p /opt/flink/plugins/s3-fs-hadoop && \
    cp /opt/flink/opt/flink-s3-fs-hadoop-${FLINK_VERSION}.jar /opt/flink/plugins/s3-fs-hadoop/

USER flink
