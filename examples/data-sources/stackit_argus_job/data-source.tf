resource "stackit_argus_instance" "example" {
  name       = "example"
  project_id = stackit_project.example.id
  plan       = "Monitoring-Medium-EU01"
}

data "stackit_argus_job" "example" {
  project_id        = stackit_project.example.id
  argus_instance_id = stackit_argus_instance.example.id
}