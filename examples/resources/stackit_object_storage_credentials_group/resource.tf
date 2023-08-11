resource "stackit_object_storage_credentials_group" "example" {
  project_id = stackit_object_storage_project.example.id
  name       = "example"
}

resource "stackit_object_storage_credential" "example" {
  project_id           = stackit_object_storage_project.example.id
  credentials_group_id = stackit_object_storage_credentials_group.example.id
}
