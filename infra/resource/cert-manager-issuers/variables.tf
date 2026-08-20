variable "vault_url" {
  description = "Vault URL — the PKI cert endpoints read here are unauthenticated by design"
  type        = string
}

variable "vault_pki_mount" {
  description = <<-EOT
    Vault PKI mount whose ca_chain is the trust anchor for every edge
    certificate. Sourcing the chain here rather than lifting ca.crt out of one
    service's leaf Secret keeps the bundle tied to the authority that issues it.
  EOT
  type        = string
  default     = "pki_int"
}

variable "trust_namespace" {
  description = "trust-manager's source namespace — Bundle sources are only read from here"
  type        = string
  default     = "cert-manager"
}

variable "bundle_name" {
  description = "Bundle name, which is also the name of the ConfigMap it writes into every namespace"
  type        = string
  default     = "internal-ca"
}

variable "tekton_bundle_name" {
  description = <<-EOT
    Tekton mounts a ConfigMap of this name into every step and points SSL_CERT_DIR at it. That
    overrides the step image's own trust store, so an image that keeps its roots elsewhere — kaniko
    uses /kaniko/ssl/certs — ends up trusting nothing at all unless this bundle exists. Same contents
    and key as the general bundle; only the name differs, because Tekton hardcodes it.
  EOT
  type        = string
  default     = "config-trusted-cabundle"
}

variable "bundle_key" {
  description = "Key inside the distributed ConfigMap. Holds the public roots plus the internal chain, so consumers point their trust store at it directly instead of concatenating"
  type        = string
  default     = "ca-bundle.crt"
}
