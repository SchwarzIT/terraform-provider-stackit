data "stackit_argus_instance" "example" {
  id         = "argus-instance-id"
  project_id = stackit_project.example.id
}