generate "provider_vault" {
  path      = "provider_vault.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<-EOF
    provider "vault" {
      address          = var.vault_url
      skip_child_token = true
      auth_login {
        path = "auth/kubernetes/login"
        parameters = {
          role = var.vault_auth_role
          jwt  = var.vault_k8s_jwt
        }
      }
    }
  EOF
}
