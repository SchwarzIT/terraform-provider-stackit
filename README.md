# Terraform Provider for STACKIT

[![Go Report Card](https://goreportcard.com/badge/github.com/SchwarzIT/terraform-provider-stackit)](https://goreportcard.com/report/github.com/SchwarzIT/terraform-provider-stackit) <!--workflow-badge-->[![GitHub Workflow Status](https://img.shields.io/badge/Acceptance%20Tests-46%20passed%2C%2032%20failed-orange)](https://github.com/SchwarzIT/terraform-provider-stackit/actions/workflows/acceptance_test.yml)<!--revision-d21dd9f5-9ca9-45e9-a767-298330b5e643--><!--workflow-badge--><br />[![GitHub release (latest by date)](https://img.shields.io/github/v/release/SchwarzIT/terraform-provider-stackit)](https://registry.terraform.io/providers/SchwarzIT/stackit/latest/docs) [![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens) [![License](https://img.shields.io/badge/License-Apache_2.0-lightgray.svg)](https://opensource.org/licenses/Apache-2.0) 

This provider is built and maintained by the STACKIT community in Schwarz IT and is not an official STACKIT provider

> The provider is built using Terraform's plugin framework, therefore we recommend using [Terraform v1.x](https://www.terraform.io/downloads) which supports Protocol v6

## Usage Example

```hcl

terraform {
  required_providers {
    stackit = {
      source  = "SchwarzIT/stackit"
      version = ">= 1.15"
    }
  }
}

# Configure the STACKIT Provider
provider "stackit" {
  service_account_email = var.service_account_email
  service_account_token = var.service_account_token
}

```

## External Links

* [Provider Documentation](https://registry.terraform.io/providers/SchwarzIT/stackit/latest/docs)
* [STACKIT Community Go Client](https://github.com/SchwarzIT/community-stackit-go-client)
