path "pki*" { capabilities = ["read", "list"] }
path "pki/roles/*" { capabilities = ["create", "update"] }
path "pki/sign/*" { capabilities = ["create", "update"] }
path "pki/issue/*" { capabilities = ["create"] }
