path "pki*" { capabilities = ["list", "read"] }
path "pki/roles/*" { capabilities = ["create", "update"] }
path "pki/sign/*" { capabilities = ["create", "update"] }
path "pki/issue/*" { capabilities = ["create"] }
