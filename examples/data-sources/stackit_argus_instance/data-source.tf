data "stackit_argus_instance" "example" {
  name       = "example"
  project_id = stackit_project.example.id
}