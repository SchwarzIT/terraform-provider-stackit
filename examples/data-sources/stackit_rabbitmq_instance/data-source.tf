resource "stackit_rabbitmq_instance" "example" {
  name       = "example"
  project_id = "example"
  version    = "3.7"
  plan       = "stackit-rabbitmq-single-small"
}
