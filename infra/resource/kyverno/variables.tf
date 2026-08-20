variable "vault_url" {
  description = "Vault URL"
  type        = string
}

variable "vault_auth_role" {
  description = "Vault Kubernetes auth role this unit authenticates as"
  type        = string
}

variable "vault_k8s_jwt" {
  description = "Short-lived ServiceAccount token used for the Vault Kubernetes auth login"
  type        = string
  sensitive   = true
}

variable "namespace" {
  description = "Kyverno namespace — registry credentials must live here for Kyverno to read them"
  type        = string
}

variable "image_pull_secret" {
  description = "Name of the dockerconfigjson Secret Kyverno uses to pull signatures and attestations"
  type        = string
}

variable "harbor_address" {
  description = "Harbor hostname as it appears in the image references Kubernetes admits"
  type        = string
}

variable "harbor_user" {
  description = "Harbor user whose credentials Kyverno reads signatures with"
  type        = string
}

variable "cosign_key_name" {
  description = "Vault Transit key that signs release images; its public half is embedded in the policies"
  type        = string
}

variable "signed_namespaces" {
  description = "Namespaces whose pods are checked for image signatures"
  type        = list(string)
}

variable "image_glob" {
  description = "Image reference glob the policies apply to. Images outside it are ignored, which is what keeps sidecars and third-party images out of scope"
  type        = string
}

variable "validation_action" {
  description = "Audit reports violations without blocking admission; Deny blocks. Start at Audit until every running image is signed"
  type        = string
}

variable "rekor_url" {
  description = <<-EOT
    Placeholder transparency-log URL. There is no Rekor here and insecureIgnoreTlog is on, so nothing
    dials it — but Kyverno's opts builder returns "rekor URL must be provided" whenever a ctlog block
    is present without one, and it evaluates that before it ever reads insecureIgnoreTlog. The pub-key
    fields are filled with the release key for the same reason: left empty they trigger a fetch from
    the public sigstore TUF, which this cluster cannot reach.
  EOT
  type        = string
}

variable "admission_enabled" {
  description = <<-EOT
    Whether the policies run at admission. Keep false while auditing: an ImageValidatingPolicy
    that evaluates at admission rejects pods whose images it could not verify — it requires its
    own verification-outcome annotation — and that rejection happens regardless of
    validation_action, so Audit alone does not make admission safe.
  EOT
  type        = bool
}

variable "provenance_predicate_type" {
  description = "in-toto predicate type Tekton Chains attaches for build provenance"
  type        = string
}

variable "builder_id" {
  description = "builder.id Tekton Chains stamps into every provenance predicate. Asserting it is what distinguishes an image our CI built from one signed by the same key out of band"
  type        = string
}

variable "provenance_policy_enabled" {
  description = <<-EOT
    Kyverno 1.18 cannot verify an attestation that has no Rekor entry: its attestation path treats
    cosign's bundleVerified=false as fatal, while the signature path exempts it when
    insecureIgnoreTlog is set (pkg/image/verifiers/ivpol/cosign/verifier.go). The attestation itself
    is sound — cosign verify-attestation accepts it against this same key. Turn this back on once a
    Rekor instance exists or Kyverno grows the matching exemption; leaving the policy running only
    produces a failure that can never clear, which teaches readers to ignore the reports.
  EOT
  type        = bool
}
