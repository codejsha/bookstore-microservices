provider "kubernetes" {
  config_path = "~/.kube/config"
}

resource "tls_private_key" "proxy_key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_cert_request" "proxy_csr" {
  private_key_pem = tls_private_key.proxy_key.private_key_pem
  subject {
    common_name  = "system:node:example.com"
    country      = "KR"
    province     = "Seoul"
    locality     = "Gangnam"
    organization = "system:nodes"
  }
  dns_names = [
    "localhost",
    "*.example.com",
  ]
  ip_addresses = [
    "127.0.0.1",
    "192.168.0.50"
  ]
}

resource "kubernetes_certificate_signing_request_v1" "proxy_kube_csr" {
  metadata {
    name = "proxy-kube-csr"
  }
  spec {
    signer_name = "kubernetes.io/kubelet-serving"
    request     = tls_cert_request.proxy_csr.cert_request_pem
    usages = [
      "digital signature",
      "key encipherment",
      "server auth",
    ]
  }
  auto_approve = true
}

resource "local_file" "proxy_key" {
  filename = "proxy.key"
  content  = tls_private_key.proxy_key.private_key_pem
}

resource "local_file" "proxy_crt" {
  filename = "proxy.crt"
  content  = kubernetes_certificate_signing_request_v1.proxy_kube_csr.certificate
}
