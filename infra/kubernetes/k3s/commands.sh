######################################################################

scp root@workstation.internal:/var/lib/rancher/k3s/server/tls/server-ca.crt ~/.kube/k3s/server-ca.crt
scp root@workstation.internal:/var/lib/rancher/k3s/server/tls/client-admin.crt ~/.kube/k3s/client-admin.crt
scp root@workstation.internal:/var/lib/rancher/k3s/server/tls/client-admin.key ~/.kube/k3s/client-admin.key

######################################################################

scp root@workstation.internal:/etc/rancher/k3s/k3s.yaml ~/.kube/k3s/k3s.yaml
yq -i '
  .clusters[0].name = "k3s" |
  .clusters[0].cluster.server = "https://workstation.internal:6443" |
  .users[0].name = "k3s" |
  .contexts[0].name = "k3s" |
  .contexts[0].context.cluster = "k3s" |
  .contexts[0].context.user = "k3s" |
  .current-context = "k3s"
' ~/.kube/k3s/k3s.yaml
yq -i eval-all '
  select(fileIndex == 0) as $host |
  select(fileIndex == 1) as $k3s |
  $host |
  .clusters = ([.clusters[] | select(.name != "k3s")] + $k3s.clusters) |
  .users = ([.users[] | select(.name != "k3s")] + $k3s.users) |
  .contexts = ([.contexts[] | select(.name != "k3s")] + $k3s.contexts)
' ~/.kube/config ~/.kube/k3s/k3s.yaml
kubectl config use-context k3s
