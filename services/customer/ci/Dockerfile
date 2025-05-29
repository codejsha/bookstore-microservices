FROM golang:1.23 AS build

WORKDIR /service
COPY . /service
RUN make build

FROM alpine:3.21 AS runtime

ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}

WORKDIR /app
COPY --from=build /service/bin/${SERVICE_NAME}-app /app/${SERVICE_NAME}-app

EXPOSE 8080
EXPOSE 9090
ENTRYPOINT ["/bin/sh", "-c", "/app/${SERVICE_NAME}-app run"]
