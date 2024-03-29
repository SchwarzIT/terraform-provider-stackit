---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "stackit_object_storage_credential Data Source - stackit"
subcategory: ""
description: |-
  Data source for Object Storage credentials
  
  -> Environment supportTo set a custom API base URL, set STACKITOBJECTSTORAGE_BASEURL environment variable
---

# stackit_object_storage_credential (Data Source)

Data source for Object Storage credentials

<br />

-> __Environment support__<small>To set a custom API base URL, set <code>STACKIT_OBJECT_STORAGE_BASEURL</code> environment variable </small>

## Example Usage

```terraform
resource "stackit_object_storage_credential" "example" {
  object_storage_project_id = stackit_object_storage_project.example.id
}

data "stackit_object_storage_credentials_group" "example" {
  depends_on = [stackit_object_storage_credential.example]
  project_id = var.project_id
  name       = "default"
}

data "stackit_object_storage_credential" "ex1" {
  project_id           = var.project_id
  credentials_group_id = data.stackit_object_storage_credentials_group.example.id
  id                   = stackit_object_storage_credential.example.id
}

data "stackit_object_storage_credential" "ex2" {
  project_id           = var.project_id
  credentials_group_id = data.stackit_object_storage_credentials_group.example.id
  display_name         = stackit_object_storage_credential.example.display_name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `credentials_group_id` (String) the credentials group ID
- `project_id` (String) The ID returned from `stackit_object_storage_project`

### Optional

- `display_name` (String) the credential's display name in the portal
- `id` (String) the credential ID
- `object_storage_project_id` (String, Deprecated) The ID returned from `stackit_object_storage_project`

### Read-Only

- `expiry` (String)


