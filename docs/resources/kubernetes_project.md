---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "stackit_kubernetes_project Resource - stackit"
subcategory: ""
description: |-
  This resource enables STACKIT Kubernetes Engine (SKE) in a project
---

# stackit_kubernetes_project (Resource)

This resource enables STACKIT Kubernetes Engine (SKE) in a project

## Example Usage

```terraform
resource "stackit_kubernetes_project" "example" {
  project_id = "example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `project_id` (String) the project ID that SKE will be enabled in

### Read-Only

- `id` (String) kubernetes project ID


