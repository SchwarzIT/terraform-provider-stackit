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
