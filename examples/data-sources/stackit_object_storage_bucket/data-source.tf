resource "stackit_object_storage_project" "example" {
  project_id = "example"
}

resource "stackit_object_storage_bucket" "example" {
  object_storage_project_id = stackit_object_storage_project.example.id
  name                      = "example"
}

data "stackit_object_storage_bucket" "example" {
  depends_on                = [stackit_object_storage_bucket.example]
  object_storage_project_id = stackit_object_storage_project.example.id
  name                      = "example"
}
