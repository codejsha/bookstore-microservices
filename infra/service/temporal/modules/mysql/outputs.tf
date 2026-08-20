output "connect_addr" {
  description = "host:port for Temporal's SQL store (read-write Service)"
  value       = "${var.cluster_name}.${var.namespace}.svc.cluster.local:3306"
}

output "db_secret_name" {
  description = "Secret name for the Temporal chart's existingSecret"
  value       = kubernetes_secret_v1.db.metadata[0].name
}
