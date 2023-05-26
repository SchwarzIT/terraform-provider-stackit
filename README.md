# Terraform Provider for STACKIT
<!--summary-image-->
<img src="https://hcti.io/v1/image/017295fd-4691-49cf-bd42-e5c53f06ae97" width="250" align="right" />
<!--revision-d2cc54b6-9902-40bd-90e7-8269ff2fbd28--><!--summary-image-->

[![Go Report Card](https://goreportcard.com/badge/github.com/SchwarzIT/terraform-provider-stackit)](https://goreportcard.com/report/github.com/SchwarzIT/terraform-provider-stackit) <!--workflow-badge-->[![GitHub Workflow Status](https://img.shields.io/badge/Acceptance%20Tests-45%20passed%2C%206%20failed-green)](https://github.com/SchwarzIT/terraform-provider-stackit/actions/workflows/acceptance_test.yml)<!--revision-7f74504c-5403-4a13-bb9c-7943dcef4171--><!--workflow-badge--><br />[![GitHub release (latest by date)](https://img.shields.io/github/v/release/SchwarzIT/terraform-provider-stackit)](https://registry.terraform.io/providers/SchwarzIT/stackit/latest/docs) [![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens) [![License](https://img.shields.io/badge/License-Apache_2.0-lightgray.svg)](https://opensource.org/licenses/Apache-2.0)

The STACKIT provider is a project developed and maintained by the STACKIT community within Schwarz IT. Please note that it is not an official provider endorsed or maintained by STACKIT.

📖 [Provider Documentation](https://registry.terraform.io/providers/SchwarzIT/stackit/latest/docs)

🚀 [STACKIT Community Go Client](https://github.com/SchwarzIT/community-stackit-go-client)

&nbsp;

## Reporting Issues

If you encounter any issues or have suggestions for improvements, please open an issue in the repository. We appreciate your feedback and will do our best to address any problems as soon as possible.

The provider is built and maintained by the STACKIT community in Schwarz IT and is not an official STACKIT provider

&nbsp;

## Usage Example

```hcl

terraform {
  required_providers {
    stackit = {
      source  = "SchwarzIT/stackit"
      version = "~> 1.19.0"
    }
  }
}

# Configure the STACKIT Provider
provider "stackit" {
  service_account_key_path = var.service_account_key_path
  private_key_path         = var.private_key_path
}

```

For further authentication methods, please refer to our [Provider Documentation](https://registry.terraform.io/providers/SchwarzIT/stackit/latest/docs)
