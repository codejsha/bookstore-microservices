#!/bin/sh

if [ -f /vault/secrets/env-keycloak ]; then
    . /vault/secrets/env-keycloak
fi

exec "/app/${SERVICE_NAME}-app" run
