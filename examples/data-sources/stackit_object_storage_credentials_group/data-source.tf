resource "stackit_object_storage_project" "example" {
  project_id = "example"
}

resource "stackit_object_storage_credentials_group" "example" {
  object_storage_project_id = stackit_object_storage_project.example.id
  name                      = "example"
}

data "stackit_object_storage_credentials_group" "example" {
  object_storage_project_id = stackit_object_storage_project.example.id
  id                        = stackit_object_storage_credentials_group.example.id
}
