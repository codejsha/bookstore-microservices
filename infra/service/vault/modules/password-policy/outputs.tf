output "policy_special" {
  description = "Password policy name for credentials that may contain special chars (admin/service accounts)."
  value       = vault_password_policy.special.name
}

output "policy_alphanumeric" {
  description = "Password policy name for DSN/CLI-safe credentials (database, S3, connection strings)."
  value       = vault_password_policy.alphanumeric.name
}
