---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "stackit_mongodb_flex_instance Resource - stackit"
subcategory: ""
description: |-
  Manages MongoDB Flex instances
  ~> Note: MongoDB Flex is in 'beta' stage in STACKIT
---

# stackit_mongodb_flex_instance (Resource)

Manages MongoDB Flex instances
		
~> **Note:** MongoDB Flex is in 'beta' stage in STACKIT

## Example Usage

```terraform
resource "stackit_mongodb_flex_instance" "example" {
  name         = "example"
  project_id   = "example"
  machine_type = "C1.1"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `machine_type` (String) The Machine Type. Available options: `T1.2`, `C1.1`, `G1.1`, `M1.1`, `C1.2`, `G1.2`, `M1.2`, `C1.3`, `G1.3`, `M1.3`, `C1.4`, `G1.4`, `M1.4`, `C1.5`, `G1.5`
- `name` (String) Specifies the instance name. Changing this value requires the resource to be recreated.
- `project_id` (String) The project ID the instance runs in. Changing this value requires the resource to be recreated.

### Optional

- `acl` (List of String) Access Control rules to whitelist IP addresses
- `backup_schedule` (String) Specifies the backup schedule (cron style)
- `labels` (Map of String) Instance Labels
- `options` (Map of String) Specifies mongodb instance options
- `replicas` (Number) Number of replicas (Default is `1`)
- `storage` (Attributes) A signle `storage` block as defined below. (see [below for nested schema](#nestedatt--storage))
- `version` (String) MongoDB version. Version `5.0` and `6.0` are supported. Changing this value requires the resource to be recreated.

### Read-Only

- `id` (String) Specifies the resource ID
- `user` (Attributes) The databse admin user (see [below for nested schema](#nestedatt--user))

<a id="nestedatt--storage"></a>
### Nested Schema for `storage`

Optional:

- `class` (String) Specifies the storage class. Available option: `premium-perf2-mongodb`
- `size` (Number) The storage size in GB. Default is `10`.


<a id="nestedatt--user"></a>
### Nested Schema for `user`

Read-Only:

- `database` (String) Specifies the database the user can access
- `host` (String) Specifies the allowed user hostname
- `id` (String) Specifies the user id
- `password` (String, Sensitive) Specifies the user's password
- `port` (Number) Specifies the port
- `roles` (List of String) Specifies the roles assigned to the user
- `uri` (String, Sensitive) Specifies connection URI
- `username` (String) Specifies the user's username


