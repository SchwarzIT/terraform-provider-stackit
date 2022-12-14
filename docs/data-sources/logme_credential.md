---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "stackit_logme_credential Data Source - stackit"
subcategory: ""
description: |-
  Manages LogMe credentials
---

# stackit_logme_credential (Data Source)

Manages LogMe credentials

## Example Usage

```terraform
resource "stackit_logme_instance" "example" {
  name       = "example"
  project_id = "example"
}

resource "stackit_logme_credential" "example" {
  project_id  = "example"
  instance_id = stackit_logme_instance.example.id
}

data "stackit_logme_credential" "example" {
  id          = stackit_logme_credential.example.id
  project_id  = "example"
  instance_id = stackit_logme_instance.example.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Specifies the resource ID
- `instance_id` (String) Instance ID the credential belongs to
- `project_id` (String) Project ID the credential belongs to

### Read-Only

- `host` (String) Credential host
- `hosts` (List of String) Credential hosts
- `password` (String) Credential password
- `port` (Number) Credential port
- `route_service_url` (String) Credential route_service_url
- `syslog_drain_url` (String) Credential syslog_drain_url
- `uri` (String) The instance URI
- `username` (String) Credential username


