#!/bin/sh

curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION="v1.36.0+k3s1" INSTALL_K3S_SKIP_SELINUX_RPM="true" K3S_KUBECONFIG_MODE="644" INSTALL_K3S_EXEC="server --cluster-init --flannel-backend=none --cluster-cidr=10.42.0.0/16 --service-cidr=10.43.0.0/16 --disable-network-policy --disable=traefik --kubelet-arg=max-pods=250" sh -
