output "approle_credentials" {
  description = "Per-role { role_id, secret_id } map for vault approle. Keys match infra/service/vault/modules/approle/main.tf:local.all_policies."
  value       = module.approle.credentials
  sensitive   = true
}

output "seaweedfs_s3_credentials" {
  description = "Pre-generated S3 admin creds (access_key_id, secret_access_key)"
  value = {
    access_key_id     = random_password.seaweedfs_s3_access_key.result
    secret_access_key = random_password.seaweedfs_s3_secret_key.result
  }
  sensitive = true
}
