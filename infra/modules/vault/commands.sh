######################################################################

scp example-ca.crt root@workstation.internal:/etc/pki/ca-trust/source/anchors/example-ca.crt
scp example-int-ca.crt root@workstation.internal:/etc/pki/ca-trust/source/anchors/example-int-ca.crt
ssh root@workstation.internal update-ca-trust
