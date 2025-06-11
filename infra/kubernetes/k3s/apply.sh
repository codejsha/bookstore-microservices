#!/bin/sh

# curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION="v1.32.5+k3s1" sh -
# curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION="v1.32.5+k3s1" INSTALL_K3S_EXEC="server --cluster-cidr=10.42.0.0/16 --service-cidr=10.43.0.0/16 --disable=traefik" sh -
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION="v1.32.5+k3s1" K3S_KUBECONFIG_MODE="644" INSTALL_K3S_EXEC="server --cluster-init --flannel-backend=none --cluster-cidr=192.168.0.0/16 --disable-network-policy --disable=traefik" sh -
