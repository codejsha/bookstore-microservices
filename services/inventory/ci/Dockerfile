FROM alpine:3.21 AS runtime

ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}
ENV TZ="Asia/Seoul"

WORKDIR /app
COPY ./bin/${SERVICE_NAME}-app .
COPY ./ci/entrypoint.sh .
RUN chmod +x entrypoint.sh

EXPOSE 8080
EXPOSE 9090
ENTRYPOINT ["/app/entrypoint.sh"]
