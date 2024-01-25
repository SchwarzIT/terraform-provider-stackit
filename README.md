# Terraform Provider for STACKIT
<!--summary-image-->
<img src="https://hcti.io/v1/image/2b3ff94c-f36c-4e9c-9566-015ab654f1b3" width="250" align="right" />
<!--revision-2c63b379-9734-46be-999c-b951550c3811--><!--summary-image-->

[![Go Report Card](https://goreportcard.com/badge/github.com/SchwarzIT/terraform-provider-stackit)](https://goreportcard.com/report/github.com/SchwarzIT/terraform-provider-stackit) <!--workflow-badge-->[![GitHub Workflow Status](https://img.shields.io/badge/Acceptance%20Tests-33%20passed%2C%2024%20failed-orange)](https://github.com/SchwarzIT/terraform-provider-stackit/actions/workflows/acceptance_test.yml)<!--revision-0ee2e4d1-7acb-4bf5-9c7a-9c31b4293387--><!--workflow-badge--><br />[![GitHub release (latest by date)](https://img.shields.io/github/v/release/SchwarzIT/terraform-provider-stackit)](https://registry.terraform.io/providers/SchwarzIT/stackit/latest/docs) [![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens) [![License](https://img.shields.io/badge/License-Apache_2.0-lightgray.svg)](https://opensource.org/licenses/Apache-2.0)

The STACKIT provider is a project developed and maintained by the STACKIT community within Schwarz IT.

**⚠️ Important:** This project will continue be maintained as long as internal teams are using it.

> Note that the official STACKIT provider can be found [here](https://github.com/stackitcloud/terraform-provider-stackit).<br />To migrate to the official provider from the community provider, [refer to the migration notes](https://github.com/stackitcloud/terraform-provider-stackit/blob/main/MIGRATION.md).

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
      version = "~> 1.27.0"
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
