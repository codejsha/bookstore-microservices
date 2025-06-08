terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
  }
}

resource "vault_mount" "pki" {
  type                  = "pki"
  path                  = "pki"
  description           = "PKI Root CA"
  max_lease_ttl_seconds = 315360000 # 10 years
}

resource "vault_pki_secret_backend_root_cert" "pki" {
  backend     = vault_mount.pki.path
  type        = "internal"
  common_name = "example corp root authority"
  ttl         = "87600h" # 10 years
}

resource "vault_pki_secret_backend_config_urls" "pki" {
  backend = vault_mount.pki.path
  issuing_certificates = ["${var.vault_url}/v1/pki/ca"]
  crl_distribution_points = ["${var.vault_url}/v1/pki/crl"]
}

resource "vault_policy" "pki" {
  name = "pki"
  policy = file("${path.module}/pki-policy.hcl")
}

resource "local_file" "pki" {
  filename = "example-ca.crt"
  content  = vault_pki_secret_backend_root_cert.pki.certificate
}

resource "vault_mount" "pki_int" {
  type        = "pki"
  path        = "pki_int"
  description = "PKI Intermediate CA"
}

resource "vault_pki_secret_backend_intermediate_cert_request" "pki_int" {
  backend     = vault_mount.pki_int.path
  type        = "exported"
  common_name = "example corp intermediate authority"
}

resource "vault_pki_secret_backend_root_sign_intermediate" "pki_int" {
  backend     = vault_mount.pki.path
  common_name = vault_pki_secret_backend_intermediate_cert_request.pki_int.common_name
  csr         = vault_pki_secret_backend_intermediate_cert_request.pki_int.csr
  format      = "pem_bundle"
  ttl         = "43800h" # 5 years
}

resource "vault_pki_secret_backend_intermediate_set_signed" "pki_int" {
  backend     = vault_mount.pki_int.path
  certificate = vault_pki_secret_backend_root_sign_intermediate.pki_int.certificate
}

resource "vault_pki_secret_backend_config_urls" "pki_int" {
  backend = vault_mount.pki_int.path
  issuing_certificates = ["${var.vault_url}/v1/pki_int/ca"]
  crl_distribution_points = ["${var.vault_url}/v1/pki_int/crl"]
}

resource "vault_policy" "pki_int" {
  name = "pki_int"
  policy = file("${path.module}/pki-int-policy.hcl")
}

resource "local_file" "pki_int" {
  filename = "example-int-ca.crt"
  content  = vault_pki_secret_backend_intermediate_set_signed.pki_int.certificate
}
