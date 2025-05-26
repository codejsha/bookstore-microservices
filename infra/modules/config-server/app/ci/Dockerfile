FROM openjdk:21 AS build

ENV SERVICE_NAME config-server
WORKDIR /service
COPY . /service
RUN gradle build

FROM openjdk:21 AS runtime

WORKDIR /app
COPY --from=build /service/build/libs/${SERVICE_NAME}.jar /app/${SERVICE_NAME}.jar

EXPOSE 8888
ENTRYPOINT ["java", "-jar", "/app/${SERVICE_NAME}.jar"]
CMD ["java"]
