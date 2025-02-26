path "pki_int*" { capabilities = ["read", "list"] }
path "pki_int/roles/*" { capabilities = ["create", "update"] }
path "pki_int/sign/*" { capabilities = ["create", "update"] }
path "pki_int/issue/*" { capabilities = ["create"] }
