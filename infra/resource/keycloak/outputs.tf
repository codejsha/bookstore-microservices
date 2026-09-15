output "bookstore_admin_initial_passwords" {
  description = "One-time initial passwords of the bookstore realm admin console accounts"
  value       = module.bookstore_admins.initial_passwords
  sensitive   = true
}
