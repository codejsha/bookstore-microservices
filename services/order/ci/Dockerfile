FROM gradle:8.13-jdk21 AS build

WORKDIR /service
COPY . /service
RUN gradle build

FROM openjdk:21 AS runtime

ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}

WORKDIR /app
COPY --from=build /service/build/libs/${SERVICE_NAME}-app.jar /app/${SERVICE_NAME}-app.jar
COPY ./ci/entrypoint.sh .
RUN chmod +x entrypoint.sh

EXPOSE 8080
EXPOSE 9090
ENTRYPOINT ["/app/entrypoint.sh"]
