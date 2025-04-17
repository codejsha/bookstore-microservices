#!/bin/bash
trap 'echo "${BASH_SOURCE[0]}: line ${LINENO}: status ${?}: user ${USER}: func ${FUNCNAME[0]}"' ERR
set -o errexit
set -o errtrace

USERNAME=""
TOKEN=""

cd catalog
git init
git checkout -b main
git add .
git commit -m "feat: initial commit"
git push -u http://${USERNAME}:${TOKEN}@git.example.com/example-corp/tekton-custom-catalog.git main
