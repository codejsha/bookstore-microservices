output "connect_addr" {
  description = "host:port for Temporal's SQL stores (read-write Service)"
  value       = "${var.cluster_name}-rw.${var.namespace}.svc.cluster.local:5432"
}

output "db_secret_name" {
  description = "CNPG superuser Secret (key `password`) for the Temporal chart's existingSecret"
  value       = "${var.cluster_name}-superuser"
}

output "db_user" {
  description = "Superuser Temporal connects as, so it can create databases and manage schemas"
  value       = "postgres"
}
