generate "provider_kubernetes" {
  path      = "provider_kubernetes.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<-EOF
    provider "kubernetes" {
      config_path = "~/.kube/config"
    }
  EOF
}
