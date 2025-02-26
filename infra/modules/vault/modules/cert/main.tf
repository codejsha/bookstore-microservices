resource "tls_private_key" "vault_key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_cert_request" "vault_csr" {
  private_key_pem = tls_private_key.vault_key.private_key_pem
  subject {
    common_name  = "system:node:vault.example.com"
    country      = "KR"
    province     = "Seoul"
    locality     = "Gangnam"
    organization = "system:nodes"
  }
  dns_names = [
    "localhost",
    "vault.example.com",
    "vault.vault.svc",
    "*.vault.vault.svc",
    "vault.vault.svc.cluster.local",
    "*.vault.vault.svc.cluster.local",
    "vault-internal.vault.svc",
    "*.vault-internal.vault.svc",
    "vault-internal.vault.svc.cluster.local",
    "*.vault-internal.vault.svc.cluster.local"
  ]
  ip_addresses = [
    "127.0.0.1"
  ]
}

resource "kubernetes_certificate_signing_request_v1" "vault_kube_csr" {
  metadata {
    name = "vault-kube-csr"
  }
  spec {
    signer_name = "kubernetes.io/kubelet-serving"
    request     = tls_cert_request.vault_csr.cert_request_pem
    usages = [
      "digital signature",
      "key encipherment",
      "server auth",
    ]
  }
  auto_approve = true
}

resource "kubernetes_secret" "vault_ha_tls" {
  type = "kubernetes.io/tls"
  metadata {
    name      = "vault-ha-tls"
    namespace = var.namespace
  }
  data = {
    "tls.key"           = local_file.vault_key.content
    "tls.crt"           = local_file.vault_crt.content
    "kubernetes-ca.crt" = var.kube_ca_crt
  }
}

resource "local_file" "vault_key" {
  filename = "vault.key"
  content  = tls_private_key.vault_key.private_key_pem
}

resource "local_file" "vault_crt" {
  filename = "vault.crt"
  content  = kubernetes_certificate_signing_request_v1.vault_kube_csr.certificate
}
