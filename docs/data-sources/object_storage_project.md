---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "stackit_object_storage_project Data Source - stackit"
subcategory: ""
description: |-
  Data source for Object Storage project
---

# stackit_object_storage_project (Data Source)

Data source for Object Storage project

## Example Usage

```terraform
resource "stackit_object_storage_project" "example" {
  project_id = "example"
}

data "stackit_object_storage_project" "example" {
  depends_on = [stackit_object_storage_project.example]
  project_id = "example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `project_id` (String) The project ID in which Object Storage is enabled

### Read-Only

- `id` (String) Specifies the resource ID


