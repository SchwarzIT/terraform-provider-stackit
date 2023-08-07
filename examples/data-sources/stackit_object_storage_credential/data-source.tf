resource "stackit_object_storage_project" "example" {
  project_id = "example"
}

resource "stackit_object_storage_credential" "example" {
  object_storage_project_id = stackit_object_storage_project.example.id
}

data "stackit_object_storage_credentials_group" "example" {
  depends_on                = [stackit_object_storage_credential.example]
  object_storage_project_id = stackit_object_storage_project.example.id
  name                      = "default"
}

data "stackit_object_storage_credential" "ex1" {
  object_storage_project_id = stackit_object_storage_project.example.id
  credentials_group_id      = data.stackit_object_storage_credentials_group.example.id
  id                        = stackit_object_storage_credential.example.id
}

data "stackit_object_storage_credential" "ex2" {
  object_storage_project_id = stackit_object_storage_project.example.id
  credentials_group_id      = data.stackit_object_storage_credentials_group.example.id
  display_name              = stackit_object_storage_credential.example.display_name
}
