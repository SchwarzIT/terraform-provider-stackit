# Terraform Provider for STACKIT
<!--summary-image-->
<img src="https://hcti.io/v1/image/d21c464d-a340-42e1-be75-cdd9626ccad7" width="250" align="right" />
<!--revision-9273e2b6-a527-4411-b05e-5a953b1d0ff2--><!--summary-image-->

[![Go Report Card](https://goreportcard.com/badge/github.com/SchwarzIT/terraform-provider-stackit)](https://goreportcard.com/report/github.com/SchwarzIT/terraform-provider-stackit) <!--workflow-badge-->[![GitHub Workflow Status](https://img.shields.io/badge/Acceptance%20Tests-70%20passed%2C%208%20failed-green)](https://github.com/SchwarzIT/terraform-provider-stackit/actions/workflows/acceptance_test.yml)<!--revision-ccb09996-27e1-4488-a119-2fa310c4d5a8--><!--workflow-badge--><br />[![GitHub release (latest by date)](https://img.shields.io/github/v/release/SchwarzIT/terraform-provider-stackit)](https://registry.terraform.io/providers/SchwarzIT/stackit/latest/docs) [![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens) [![License](https://img.shields.io/badge/License-Apache_2.0-lightgray.svg)](https://opensource.org/licenses/Apache-2.0)

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
