data "stackit_object_storage_credentials_group" "example" {
  id         = "..."
  project_id = stackit_project.example.id
}
