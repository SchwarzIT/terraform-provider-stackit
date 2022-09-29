# Terraform Provider for STACKIT

[![Go Report Card](https://goreportcard.com/badge/github.com/SchwarzIT/terraform-provider-stackit)](https://goreportcard.com/report/github.com/SchwarzIT/terraform-provider-stackit) [![Build Status](https://dev.azure.com/schwarzit/schwarzit.odj.core/_apis/build/status/Stackit/Stackit%20E2E%20Test?branchName=main&label=E2E%20Tests)](https://dev.azure.com/schwarzit/schwarzit.odj.core/_build/latest?definitionId=17957&branchName=main) ![GitHub release (latest by date)](https://img.shields.io/github/v/release/SchwarzIT/terraform-provider-stackit) [![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens) [![License](https://img.shields.io/badge/License-Apache_2.0-lightgray.svg)](https://opensource.org/licenses/Apache-2.0) 

This provider is built and maintained by the STACKIT community in Schwarz IT and is not an official STACKIT provider

> **_NOTE:_** The provider is built using Terraform's plugin framework, therefore we recommend using [Terraform CLI v1.x](https://www.terraform.io/downloads) which supports Protocol v6

* [Terraform Website](https://www.terraform.io)
* [STACKIT Community Go Client](https://github.com/SchwarzIT/community-stackit-go-client)
* [Provider Documentation](https://registry.terraform.io/providers/SchwarzIT/stackit/latest/docs)

## Usage Example

```hcl

terraform {
  required_providers {
    stackit = {
      source  = "SchwarzIT/stackit"
      version = "=0.1.4"
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
