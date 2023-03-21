# Terraform Provider for STACKIT
<!--summary-image-->
<img src="https://hcti.io/v1/image/d9c2fc56-16ac-4ace-8936-e307786fb46a" width="250" align="right" />
<!--revision-8f30c50a-0470-4521-a7df-1c179a26ae75--><!--summary-image-->

[![Go Report Card](https://goreportcard.com/badge/github.com/SchwarzIT/terraform-provider-stackit)](https://goreportcard.com/report/github.com/SchwarzIT/terraform-provider-stackit) <!--workflow-badge-->[![GitHub Workflow Status](https://img.shields.io/badge/Acceptance%20Tests-74%20passed%2C%204%20failed-success)](https://github.com/SchwarzIT/terraform-provider-stackit/actions/workflows/acceptance_test.yml)<!--revision-d34b6dcf-50f7-40d4-be8f-7f77827a4f71--><!--workflow-badge--><br />[![GitHub release (latest by date)](https://img.shields.io/github/v/release/SchwarzIT/terraform-provider-stackit)](https://registry.terraform.io/providers/SchwarzIT/stackit/latest/docs) [![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens) [![License](https://img.shields.io/badge/License-Apache_2.0-lightgray.svg)](https://opensource.org/licenses/Apache-2.0)

The provider is built and maintained by the STACKIT community in Schwarz IT and is not an official STACKIT provider

📖 [Provider Documentation](https://registry.terraform.io/providers/SchwarzIT/stackit/latest/docs)

🚀 [STACKIT Community Go Client](https://github.com/SchwarzIT/community-stackit-go-client)

<br />

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
