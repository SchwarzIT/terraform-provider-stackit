resource "stackit_elasticsearch_instance" "example" {
  name       = "example"
  project_id = "example"
  version    = "7"
}

resource "stackit_elasticsearch_credential" "example" {
  project_id  = "example"
  instance_id = stackit_elasticsearch_instance.example.id
}

data "stackit_elasticsearch_credential" "example" {
  id          = stackit_elasticsearch_credential.example.id
  project_id  = "example"
  instance_id = stackit_elasticsearch_instance.example.id
}
