FROM eclipse-temurin:25-jdk-alpine AS runtime

ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}
ENV TZ="Asia/Seoul"

ARG IMAGE_VERSION=1.0.0
ARG IMAGE_REVISION

LABEL org.opencontainers.image.title="${SERVICE_NAME}" \
      org.opencontainers.image.description="Order service: orders, cart, checkout and Temporal saga orchestration" \
      org.opencontainers.image.source="https://git.example.com/example-corp/${SERVICE_NAME}-source" \
      org.opencontainers.image.vendor="example-corp" \
      org.opencontainers.image.version="${IMAGE_VERSION}" \
      org.opencontainers.image.revision="${IMAGE_REVISION}"

ARG PYROSCOPE_AGENT_VERSION=0.16.0
ARG OTEL_AGENT_VERSION=2.24.0

RUN addgroup -g 65532 nonroot \
    && adduser -u 65532 -G nonroot -D nonroot

WORKDIR /app
COPY ./build/libs/${SERVICE_NAME}-app.jar .
COPY ./ci/entrypoint.sh .
ADD https://repo1.maven.org/maven2/io/pyroscope/agent/${PYROSCOPE_AGENT_VERSION}/agent-${PYROSCOPE_AGENT_VERSION}.jar /app/pyroscope.jar
ADD https://github.com/open-telemetry/opentelemetry-java-instrumentation/releases/download/v${OTEL_AGENT_VERSION}/opentelemetry-javaagent.jar /app/otel-agent.jar
RUN chmod 644 /app/pyroscope.jar /app/otel-agent.jar \
    && chmod +x entrypoint.sh

USER nonroot

EXPOSE 8080
EXPOSE 9090
ENTRYPOINT ["/app/entrypoint.sh"]
