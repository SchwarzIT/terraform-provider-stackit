resource "stackit_object_storage_project" "example" {
  project_id = "example"
}

resource "stackit_object_storage_credential" "example" {
  object_storage_project_id = stackit_object_storage_project.example.id
}

data "stackit_object_storage_credential" "ex1" {
  object_storage_project_id = stackit_object_storage_project.example.id
  id                        = stackit_object_storage_credential.example.id
}

data "stackit_object_storage_credential" "ex2" {
  object_storage_project_id = stackit_object_storage_project.example.id
  display_name              = stackit_object_storage_credential.example.display_name
}
