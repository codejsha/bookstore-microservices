######################################################################

scp root@workstation.local:/var/lib/rancher/k3s/server/tls/server-ca.crt ../../modules/vault/ca.crt
scp root@workstation.local:/var/lib/rancher/k3s/server/tls/server-ca.crt ~/.kube/k3s/server-ca.crt
