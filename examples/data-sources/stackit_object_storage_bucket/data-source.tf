data "stackit_object_storage_bucket" "example" {
  name       = "example"
  project_id = stackit_project.example.id
}
