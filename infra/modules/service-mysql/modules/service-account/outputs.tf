output "service_accounts" {
  value = [
    for service in var.mysql_services : "${service}-mysql"
  ]
}
