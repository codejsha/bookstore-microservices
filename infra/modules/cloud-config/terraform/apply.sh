#!/bin/bash
trap 'echo "${BASH_SOURCE[0]}: line ${LINENO}: status ${?}: user ${USER}: func ${FUNCNAME[0]}"' ERR
set -o errexit
set -o errtrace
set -o xtrace

/bin/cp -f ../../vault/example-int-ca.crt .

terraform init -upgrade
terraform apply -auto-approve -target module.helm
