terraform {
  required_providers {
    stackit = {
      source  = "SchwarzIT/stackit"
      version = "=1.4.0"
    }
  }
}

provider "stackit" {
  service_account_email = var.service_account_email
  service_account_token = var.service_account_token
}
