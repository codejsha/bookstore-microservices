output "bookstore_admin_initial_passwords" {
  description = "One-time initial passwords of the bookstore realm admin console accounts"
  value       = module.bookstore_admins.initial_passwords
  sensitive   = true
}

output "platform_account_initial_passwords" {
  description = "Initial passwords of the platform-infra developer accounts, which are not stored in Vault"
  value       = module.infra_realm.initial_passwords
  sensitive   = true
}
