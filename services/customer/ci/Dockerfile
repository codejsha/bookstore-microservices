FROM golang:1.23 AS build

ARG SERVICE_NAME
WORKDIR /service
COPY package /service
RUN make build

FROM alpine:3.21 AS runtime

WORKDIR /app
COPY --from=build /service/bin/${SERVICE_NAME}-app /app/${SERVICE_NAME}-app

EXPOSE 8080
EXPOSE 9090
ENTRYPOINT ["/app/${SERVICE_NAME}-app"]
CMD ["run"]
