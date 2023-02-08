resource "stackit_argus_instance" "example" {
  name       = "example"
  project_id = stackit_project.example.id
  plan       = "Monitoring-Medium-EU01"
}

resource "stackit_argus_credential" "example" {
  project_id  = stackit_project.example.id
  instance_id = stackit_argus_instance.example.id
}
