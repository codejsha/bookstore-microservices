output "cluster_names" {
  description = "Map of service name => CloudNativePG Cluster .metadata.name that was applied. Consumed by the parent module to order the Vault db-engine (which reads the CNPG superuser secret) after the clusters exist."
  value       = { for k, v in var.clusters : k => v.cluster_name }
}
