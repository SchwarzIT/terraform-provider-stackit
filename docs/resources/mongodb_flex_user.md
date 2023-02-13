---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "stackit_mongodb_flex_user Resource - stackit"
subcategory: ""
description: |-
  Manages MongoDB Flex instance users
  
  -> Environment support
  Productionhttps://api.stackit.cloud/mongodb/v1/
  QAhttps://api-qa.stackit.cloud/mongodb/v1/
  Devhttps://api-dev.stackit.cloud/mongodb/v1/
  
  By default, production is used.To set a custom URL, set an environment variable STACKITMONGODBFLEX_BASEURL
---

# stackit_mongodb_flex_user (Resource)

Manages MongoDB Flex instance users

<br />

-> __Environment support__<br /><table style='border-collapse: separate; border-spacing: 0px; margin-top:-20px; margin-left: 24px; font-size: smaller;'>
<tr><td style='width: 100px; background: #fbfcff; border: none;'>Production</td><td style='background: #fbfcff; border: none;'>https://api.stackit.cloud/mongodb/v1/</td></tr>
<tr><td style='background: #fbfcff; border: none;'>QA</td><td style='background: #fbfcff; border: none;'>https://api-qa.stackit.cloud/mongodb/v1/</td></tr>
<tr><td style='background: #fbfcff; border: none;'>Dev</td><td style='background: #fbfcff; border: none;'>https://api-dev.stackit.cloud/mongodb/v1/</td></tr>
</table><br />
<small style='margin-left: 24px; margin-top: -5px; display: inline-block;'><a href="https://registry.terraform.io/providers/SchwarzIT/stackit/latest/docs#environment">By default</a>, production is used.<br />To set a custom URL, set an environment variable <code>STACKIT_MONGODB_FLEX_BASEURL</code></small>

## Example Usage

```terraform
resource "stackit_mongodb_flex_instance" "example" {
  name         = "example"
  project_id   = var.project_id
  machine_type = "1.1"
}

resource "stackit_mongodb_flex_user" "example" {
  project_id  = var.project_id
  instance_id = stackit_mongodb_flex_instance.example.id
}

output "mongodb_username" {
  value = stackit_mongodb_flex_user.example.username
}

output "mongodb_password" {
  value     = stackit_mongodb_flex_user.example.password
  sensitive = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `instance_id` (String) the mongo db flex instance id.
- `project_id` (String) The project ID the instance runs in. Changing this value requires the resource to be recreated.

### Optional

- `database` (String) Specifies the database the user can access
- `role` (String) Specifies the role assigned to the user, either `readWrite` or `read`
- `username` (String) Specifies the user's username

### Read-Only

- `host` (String) Specifies the allowed user hostname
- `id` (String) Specifies the resource ID
- `password` (String, Sensitive) Specifies the user's password
- `port` (Number) Specifies the port
- `uri` (String, Sensitive) Specifies connection URI

