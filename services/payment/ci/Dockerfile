FROM openjdk:21 AS build

ARG SERVICE_NAME
WORKDIR /service
COPY package /service
RUN gradle build

FROM openjdk:21 AS runtime

WORKDIR /app
COPY --from=build /service/build/libs/${SERVICE_NAME}-app.jar /app/${SERVICE_NAME}-app.jar

EXPOSE 8080
EXPOSE 9090
ENTRYPOINT ["java", "-jar", "/app/${SERVICE_NAME}-app.jar"]
CMD ["java"]
