terraform {
  required_providers {
    stackit = {
      source  = "SchwarzIT/stackit"
      version = "=0.3.2"
    }
  }
}

provider "stackit" {
  service_account_id    = var.service_account_id
  service_account_token = var.service_account_token
  customer_account_id   = var.customer_account_id
}
