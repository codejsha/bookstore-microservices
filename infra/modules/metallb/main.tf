provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

resource "kubernetes_namespace" "metallb" {
  metadata {
    name = var.namespace
  }
}

resource "helm_release" "metallb" {
  namespace  = kubernetes_namespace.metallb.metadata.0.name
  name       = "metallb"
  repository = "https://metallb.github.io/metallb"
  chart      = "metallb"
  version    = "0.9.3"
}

resource "null_resource" "enable_strict_arp" {
  provisioner "local-exec" {
    command = <<EOT
    kubectl get configmap kube-proxy -n kube-system -o yaml | \
    sed -e "s/strictARP: false/strictARP: true/" | \
    kubectl apply -f - -n kube-system
    EOT
  }
  triggers = {
    helm_release_name = helm_release.metallb.metadata.0.name
  }
}

resource "kubernetes_manifest" "address_pool" {
  depends_on = [null_resource.enable_strict_arp]
  manifest = {
    apiVersion = "v1"
    kind       = "ConfigMap"
    metadata = {
      name      = "config"
      namespace = kubernetes_namespace.metallb.metadata.0.name
    }
    data = {
      config = <<EOT
      address-pools:
       - name: default
         protocol: layer2
         addresses:
          - ${var.pool_addresses}
      EOT
    }
  }
}
