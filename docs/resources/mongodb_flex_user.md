---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "stackit_mongodb_flex_user Resource - stackit"
subcategory: ""
description: |-
  Manages MongoDB Flex instance users
  
  -> Environment supportTo set a custom API base URL, set STACKITMONGODBFLEX_BASEURL environment variable
---

# stackit_mongodb_flex_user (Resource)

Manages MongoDB Flex instance users

<br />

-> __Environment support__<small>To set a custom API base URL, set <code>STACKIT_MONGODB_FLEX_BASEURL</code> environment variable </small>

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
- `roles` (List of String) Specifies the role assigned to the user, valid options are: `readWrite` or `read`
- `username` (String) Specifies the user's username

### Read-Only

- `host` (String) Specifies the allowed user hostname
- `id` (String) Specifies the resource ID
- `password` (String, Sensitive) Specifies the user's password
- `port` (Number) Specifies the port
- `uri` (String, Sensitive) Specifies connection URI


