FROM eclipse-temurin:25-jdk-alpine AS runtime

ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}
ENV TZ="Asia/Seoul"

RUN addgroup -g 65532 nonroot \
    && adduser -u 65532 -G nonroot -D nonroot

WORKDIR /app
COPY ./build/libs/${SERVICE_NAME}.jar .
COPY ./ci/entrypoint.sh .
RUN chmod +x entrypoint.sh

USER nonroot

EXPOSE 8080
EXPOSE 9090
ENTRYPOINT ["/app/entrypoint.sh"]
