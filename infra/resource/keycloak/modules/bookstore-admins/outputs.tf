output "initial_passwords" {
  description = "One-time initial passwords keyed like var.accounts; Keycloak forces a password change at first login"
  value       = { for k, v in random_password.account : k => v.result }
  sensitive   = true
}
