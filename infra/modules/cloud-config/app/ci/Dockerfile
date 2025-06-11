FROM gradle:8.13-jdk21 AS build

WORKDIR /service
COPY . /service
RUN gradle build

FROM alpine/java:21-jre AS runtime

ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}

WORKDIR /app
COPY --from=build /service/build/libs/${SERVICE_NAME}.jar .
COPY ./ci/entrypoint.sh .
RUN chmod +x entrypoint.sh

EXPOSE 8888
ENTRYPOINT ["/app/entrypoint.sh"]
