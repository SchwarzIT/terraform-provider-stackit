---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "stackit_postgres_flex_user Resource - stackit"
subcategory: ""
description: |-
  Manages Postgres Flex instance users
  
  -> Environment support
  Productionhttps://postgres-flex-service.api.eu01.stackit.cloud
  QAhttps://postgres-flex-service.api.eu01.qa.stackit.cloud
  Devhttps://postgres-flex-service.api.eu01.dev.stackit.cloud
  
  By default, production is used.To set a custom URL, set an environment variable STACKITPOSTGRESFLEX_BASEURL
---

# stackit_postgres_flex_user (Resource)

Manages Postgres Flex instance users

<br />

-> __Environment support__<br /><table style='border-collapse: separate; border-spacing: 0px; margin-top:-20px; margin-left: 24px; font-size: smaller;'>
<tr><td style='width: 100px; background: #fbfcff; border: none;'>Production</td><td style='background: #fbfcff; border: none;'>https://postgres-flex-service.api.eu01.stackit.cloud</td></tr>
<tr><td style='background: #fbfcff; border: none;'>QA</td><td style='background: #fbfcff; border: none;'>https://postgres-flex-service.api.eu01.qa.stackit.cloud</td></tr>
<tr><td style='background: #fbfcff; border: none;'>Dev</td><td style='background: #fbfcff; border: none;'>https://postgres-flex-service.api.eu01.dev.stackit.cloud</td></tr>
</table><br />
<small style='margin-left: 24px; margin-top: -5px; display: inline-block;'><a href="https://registry.terraform.io/providers/SchwarzIT/stackit/latest/docs#environment">By default</a>, production is used.<br />To set a custom URL, set an environment variable <code>STACKIT_POSTGRES_FLEX_BASEURL</code></small>

## Example Usage

```terraform
resource "stackit_postgres_flex_instance" "example" {
  name         = "example"
  project_id   = var.project_id
  machine_type = "c1.2"
  version      = "14"
}

resource "stackit_postgres_flex_user" "example" {
  project_id  = var.project_id
  instance_id = stackit_postgres_flex_instance.example.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `instance_id` (String) the postgres db flex instance id.
- `project_id` (String) The project ID the instance runs in. Changing this value requires the resource to be recreated.

### Optional

- `roles` (List of String) Specifies the roles assigned to the user, valid options are: `login`, `createdb`
- `username` (String) Specifies the username. Defaults to `psqluser`

### Read-Only

- `host` (String) Specifies the allowed user hostname
- `id` (String) Specifies the resource ID
- `password` (String, Sensitive) Specifies the user's password
- `port` (Number) Specifies the port
- `uri` (String, Sensitive) Specifies connection URI

