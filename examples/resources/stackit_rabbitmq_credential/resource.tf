resource "stackit_rabbitmq_instance" "example" {
  name       = "example"
  project_id = "example"
  version    = "3.10"
  plan       = "stackit-rabbitmq-2.4.10-single"
}

resource "stackit_rabbitmq_credential" "example" {
  project_id  = "example"
  instance_id = stackit_rabbitmq_instance.example.id
}
