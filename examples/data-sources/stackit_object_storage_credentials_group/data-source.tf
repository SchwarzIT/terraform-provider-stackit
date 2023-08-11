resource "stackit_object_storage_credentials_group" "example" {
  project_id = var.project_id
  name       = "example"
}

data "stackit_object_storage_credentials_group" "example" {
  project_id = var.project_id
  id         = stackit_object_storage_credentials_group.example.id
}
