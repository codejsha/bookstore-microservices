######################################################################

/bin/cp -f example-ca.crt /etc/pki/ca-trust/source/anchors/example-ca.crt
/bin/cp -f example-int-ca.crt /etc/pki/ca-trust/source/anchors/example-int-ca.crt
update-ca-trust
