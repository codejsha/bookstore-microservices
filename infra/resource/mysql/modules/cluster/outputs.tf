output "bootstrap_job_names" {
  description = "Names of the per-service DB/app-user bootstrap Jobs (for explicit downstream depends_on)."
  value       = { for k, j in kubernetes_job_v1.bootstrap : k => j.metadata[0].name }
}
