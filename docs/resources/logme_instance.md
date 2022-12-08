---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "stackit_logme_instance Resource - stackit"
subcategory: ""
description: |-
  Manages LogMe instances
  ~> Note: LogMe API (Part of DSA APIs) currently has issues reflecting updates & configuration correctly. Therefore, this resource is not ready for production usage.
---

# stackit_logme_instance (Resource)

Manages LogMe instances

~> **Note:** LogMe API (Part of DSA APIs) currently has issues reflecting updates & configuration correctly. Therefore, this resource is not ready for production usage.

## Example Usage

```terraform
resource "stackit_logme_instance" "example" {
  name       = "example"
  project_id = "example"
  version    = "LogMe"
  plan       = "stackit-logme-single-small-non-ssl"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Specifies the instance name. Changing this value requires the resource to be recreated. Changing this value requires the resource to be recreated.
- `project_id` (String) The project ID.

### Optional

- `acl` (List of String) Access Control rules to whitelist IP addresses
- `plan` (String) The LogMe Plan. Default is `stackit-logme-single-small-non-ssl`.
Options are: `stackit-logme-platform-logging-non-ssl`, `stackit-logme-single-small-non-ssl`, `stackit-logme-single-medium-non-ssl`, `stackit-logme-cluster-big-non-ssl`, `stackit-logme-cluster-medium-non-ssl`, `stackit-logme-cluster-small-non-ssl`
- `version` (String) LogMe version. Only Option: `LogMe`. Changing this value requires the resource to be recreated.

### Read-Only

- `cf_guid` (String) Cloud Foundry GUID
- `cf_organization_guid` (String) Cloud Foundry Organization GUID
- `cf_space_guid` (String) Cloud Foundry Space GUID
- `dashboard_url` (String) Dashboard URL
- `id` (String) Specifies the resource ID
- `plan_id` (String) The selected plan ID

