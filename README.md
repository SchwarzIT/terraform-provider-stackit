# Terraform Provider for STACKIT

This provider is built and maintained by the STACKIT community in Schwarz IT and is not an official STACKIT provider

The provider is built using Terraform's plugin framework, therefore we recommend using [Terraform CLI v1.x](https://www.terraform.io/downloads) which supports Protocol v6

* [Terraform Website](https://www.terraform.io)
* [STACKIT Community Client](https://github.com/SchwarzIT/community-stackit-go-client)

## Usage Example

```hcl

terraform {
  required_providers {
    stackit = {
      source  = "github.com/schwarzit/stackit"
      version = ">=0.0.1"
    }
  }
}

# Configure the STACKIT Provider
provider "stackit" {
  service_account_id    = var.service_account_id
  service_account_token = var.service_account_token
  customer_account_id   = var.customer_account_id
}

# create a project
resource "stackit_project" "example" {
  name                  = var.project_name
  billing_ref           = var.project_billing_ref
  owner                 = var.project_owner
  enable_kubernetes     = true
}

# create a SKE cluster
resource "stackit_kubernetes_cluster" "example" {
  name               = "my-cluster"
  project_id         = stackit_project.example.id
  kubernetes_version = "1.23.5"

  node_pools = [{
    name         = "example"
    machine_type = "c1.2"
  }]
}

```
