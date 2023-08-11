resource "stackit_object_storage_bucket" "example" {
  project_id = stackit_object_storage_project.example.id
  name       = "example"
}
