resource "stackit_object_storage_credentials_group" "example" {
  project_id = stackit_project.example.id
  name       = "example"
}
