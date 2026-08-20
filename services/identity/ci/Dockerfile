FROM alpine:3.21 AS runtime

ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}
ENV TZ="Asia/Seoul"

ARG IMAGE_VERSION=1.0.0
ARG IMAGE_REVISION

LABEL org.opencontainers.image.title="${SERVICE_NAME}" \
      org.opencontainers.image.description="Identity service: authentication and user management backed by Keycloak" \
      org.opencontainers.image.source="https://git.example.com/example-corp/${SERVICE_NAME}-source" \
      org.opencontainers.image.vendor="example-corp" \
      org.opencontainers.image.version="${IMAGE_VERSION}" \
      org.opencontainers.image.revision="${IMAGE_REVISION}"

RUN addgroup -g 65532 nonroot \
    && adduser -u 65532 -G nonroot -D nonroot

WORKDIR /app
COPY ./bin/${SERVICE_NAME}-app .
COPY ./ci/entrypoint.sh .
RUN chmod +x ${SERVICE_NAME}-app entrypoint.sh

USER nonroot

EXPOSE 8080
EXPOSE 9090
ENTRYPOINT ["/app/entrypoint.sh"]
