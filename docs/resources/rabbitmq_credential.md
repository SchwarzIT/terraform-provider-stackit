---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "stackit_rabbitmq_credential Resource - stackit"
subcategory: ""
description: |-
  Manages RabbitMQ credentials
---

# stackit_rabbitmq_credential (Resource)

Manages RabbitMQ credentials

## Example Usage

```terraform
resource "stackit_rabbitmq_instance" "example" {
  name       = "example"
  project_id = "example"
  version    = "3.7"
  plan       = "stackit-rabbitmq-single-small"
}

resource "stackit_rabbitmq_credential" "example" {
  project_id  = "example"
  instance_id = stackit_rabbitmq_instance.example.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `instance_id` (String) Instance ID the credential belongs to
- `project_id` (String) Project ID the credential belongs to

### Read-Only

- `host` (String) Credential host
- `hosts` (List of String) Credential hosts
- `id` (String) Specifies the resource ID
- `password` (String) Credential password
- `port` (Number) Credential port
- `route_service_url` (String) Credential route_service_url
- `syslog_drain_url` (String) Credential syslog_drain_url
- `uri` (String) The instance URI
- `username` (String) Credential username


