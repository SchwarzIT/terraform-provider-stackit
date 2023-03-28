# Authentication examples:

# Token flow
provider "stackit" {
  service_account_email = var.service_account_email
  service_account_token = var.service_account_token
}

# Key flow (1)
provider "stackit" {
  service_account_key_path = var.service_account_key_path
  private_key_path         = var.private_key_path
}

# Key flow (2)
provider "stackit" {
  service_account_key = var.service_account_key
  private_key         = var.private_key
}
