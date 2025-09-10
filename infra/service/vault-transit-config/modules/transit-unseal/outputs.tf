output "k8s_secret_name" {
  description = "Name of the K8s Secret carrying the autounseal token into the main Vault namespace"
  value       = kubernetes_secret_v1.autounseal.metadata[0].name
}

output "transit_mount_path" {
  description = "Transit mount path on the unsealer Vault (matches values.yaml mount_path)"
  value       = vault_mount.transit.path
}

output "transit_key_name" {
  description = "Transit key name (matches values.yaml key_name)"
  value       = vault_transit_secret_backend_key.autounseal.name
}
