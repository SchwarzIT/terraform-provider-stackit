# Terraform Provider for STACKIT

[![Go Report Card](https://goreportcard.com/badge/github.com/SchwarzIT/terraform-provider-stackit)](https://goreportcard.com/report/github.com/SchwarzIT/terraform-provider-stackit) [![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/SchwarzIT/terraform-provider-stackit/acceptance_test.yml?branch=main&label=Acceptance%20Tests)](https://github.com/SchwarzIT/terraform-provider-stackit/actions/workflows/acceptance_test.yml) [![Azure DevOps builds](https://img.shields.io/azure-devops/build/schwarzit/schwarzit.odj.core/17957?label=E2E%20Tests)](https://dev.azure.com/schwarzit/schwarzit.odj.core/_build/latest?definitionId=17957&branchName=main) <br />[![GitHub release (latest by date)](https://img.shields.io/github/v/release/SchwarzIT/terraform-provider-stackit)](https://registry.terraform.io/providers/SchwarzIT/stackit/latest/docs) [![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens) [![License](https://img.shields.io/badge/License-Apache_2.0-lightgray.svg)](https://opensource.org/licenses/Apache-2.0) 

This provider is built and maintained by the STACKIT community in Schwarz IT and is not an official STACKIT provider

> The provider is built using Terraform's plugin framework, therefore we recommend using [Terraform v1.x](https://www.terraform.io/downloads) which supports Protocol v6

## Usage Example

```hcl

terraform {
  required_providers {
    stackit = {
      source  = "SchwarzIT/stackit"
      version = "=1.8.0"
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
