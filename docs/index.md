---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "STACKIT Provider"
subcategory: ""
description: |-
  This provider is built and maintained by the STACKIT community in Schwarz IT and is not an official STACKIT provider
  ~> Note: The provider is built using Terraform's plugin framework, therefore we recommend using Terraform CLI v1.x which supports Protocol v6
---

# STACKIT Provider

This provider is built and maintained by the STACKIT community in Schwarz IT and is not an official STACKIT provider

~> **Note:** The provider is built using Terraform's plugin framework, therefore we recommend using Terraform CLI v1.x which supports Protocol v6

## Example Usage

```terraform
provider "stackit" {
  service_account_email = var.service_account_email
  service_account_token = var.service_account_token
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `service_account_email` (String) Service Account Email.<br />This attribute can also be loaded from `STACKIT_SERVICE_ACCOUNT_EMAIL` environment variable instead.
- `service_account_token` (String, Sensitive) Service Account Token.<br />This attribute can also be loaded from `STACKIT_SERVICE_ACCOUNT_TOKEN` environment variable instead.
