resource "stackit_logme_instance" "example" {
  name       = "example"
  project_id = "example"
}

resource "stackit_logme_credential" "example" {
  project_id  = "example"
  instance_id = stackit_logme_instance.example.id
}