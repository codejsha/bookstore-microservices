#!/bin/bash
trap 'echo "${bash_source[0]}: line ${lineno}: status ${?}: user ${user}: func ${funcname[0]}"' err
set -o errexit
set -o errtrace
set -o xtrace

cd ../../services/identity && make idpapi
