output "valkey_password" {
  value     = random_password.valkey.result
  sensitive = true
}

output "postgresql_username" {
  value     = var.postgresql_username
  sensitive = true
}

output "postgresql_password" {
  value     = random_password.postgresql.result
  sensitive = true
}
output "admin_password" {
  value     = random_password.admin.result
  sensitive = true
}

output "ssh_host_public_key" {
  value = trimspace(tls_private_key.ssh_host.public_key_openssh)
}

output "ssh_host_key_secret_name" {
  value = var.ssh_host_key_secret_name
}
