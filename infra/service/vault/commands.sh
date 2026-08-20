######################################################################

openssl verify -CAfile example-ca.crt example-int-ca.crt

scp example-ca.crt root@workstation.internal:/etc/pki/ca-trust/source/anchors/example-ca.crt
scp example-int-ca.crt root@workstation.internal:/etc/pki/ca-trust/source/anchors/example-int-ca.crt
ssh root@workstation.internal update-ca-trust

sudo keytool -importcert -alias example-ca -file ./example-ca.crt \
  -keystore $JAVA_HOME/lib/security/cacerts \
  -noprompt -storepass changeit

sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain example-ca.crt
sudo security add-trusted-cert -d -r trustAsRoot -k /Library/Keychains/System.keychain example-int-ca.crt
