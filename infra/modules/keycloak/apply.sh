#!/bin/bash
trap 'echo "${BASH_SOURCE[0]}: line ${LINENO}: status ${?}: user ${USER}: func ${FUNCNAME[0]}"' ERR
set -o errexit
set -o errtrace

terraform init -upgrade
# terraform import kubernetes_namespace.keycloak keycloak
# terraform import module.helm.helm_release.keycloak keycloak/keycloak

terraform apply -auto-approve -target module.realm
terraform apply -auto-approve -target module.helm || true
terraform apply -auto-approve -target module.istio
terraform apply -auto-approve -target module.cert
