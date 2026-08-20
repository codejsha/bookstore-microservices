terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
    }
  }
}

resource "vault_password_policy" "special" {
  name   = "password-special"
  policy = file("${path.module}/policies/special.hcl")
}

resource "vault_password_policy" "alphanumeric" {
  name   = "password-alphanumeric"
  policy = file("${path.module}/policies/alphanumeric.hcl")
}
