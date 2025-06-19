FROM alpine/java:21-jre AS runtime

ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}

WORKDIR /app
COPY ./build/libs/${SERVICE_NAME}-app.jar .
COPY ./ci/entrypoint.sh .
RUN chmod +x entrypoint.sh

EXPOSE 8080
EXPOSE 9090
ENTRYPOINT ["/app/entrypoint.sh"]
