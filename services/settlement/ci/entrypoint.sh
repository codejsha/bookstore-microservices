#!/bin/sh

exec java -javaagent:/app/pyroscope.jar -javaagent:/app/otel-agent.jar -jar /app/${SERVICE_NAME}-app.jar
