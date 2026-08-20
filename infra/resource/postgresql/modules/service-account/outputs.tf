output "service_accounts" {
  value = [
    for service in var.postgres_services : "${service}-postgres"
  ]
}
