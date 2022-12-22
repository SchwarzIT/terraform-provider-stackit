resource "stackit_object_storage_project" "example" {
  project_id = "example"
}

data "stackit_object_storage_project" "example" {
  depends_on = [stackit_object_storage_project.example]
  project_id = "example"
}
