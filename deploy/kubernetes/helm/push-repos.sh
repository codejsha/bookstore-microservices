#!/bin/bash
trap 'echo "${BASH_SOURCE[0]}: line ${LINENO}: status ${?}: user ${USER}: func ${FUNCNAME[0]}"' ERR
set -o errexit
set -o errtrace

USERNAME=""
TOKEN=""

repos=(
  "catalog-command"
  "catalog-query"
  "customer-command"
  "customer-query"
  "identity"
  "inventory-command"
  "inventory-query"
  "order-command"
  "order-query"
  "payment-command"
  "payment-query"
)

for repo in "${repos[@]}"; do
  cd ${repo}
  git init
  git checkout -b main
  git add .
  git commit -m "feat: initial commit"
  git push -u http://${USERNAME}:${TOKEN}@git.example.com/example-corp/${repo}-helm.git main
  cd ..
done
