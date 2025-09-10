terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.3"
    }
  }
}

module "component" {
  source = "./modules/component"
}
