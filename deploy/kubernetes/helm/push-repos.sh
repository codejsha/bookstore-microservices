#!/bin/bash
trap 'echo "${BASH_SOURCE[0]}: line ${LINENO}: status ${?}: user ${USER}: func ${FUNCNAME[0]}"' ERR
set -o errexit
set -o errtrace

USERNAME=""
TOKEN=""

service_name=("catalog" "customer" "identity" "inventory" "order" "payment")
for service in "${service_name[@]}"; do
  if [ ! -d ${service} ]; then
    exit 1
  fi

  cd ${service}
  git init
  git checkout -b main
  git add .
  git commit -m "feat: initial commit"
  git push -u http://${USERNAME}:${TOKEN}@git.example.com/example-corp/${service}-helm.git main
  cd ..
done
