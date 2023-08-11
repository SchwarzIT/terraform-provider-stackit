
resource "stackit_object_storage_bucket" "example" {
  project_id = var.project_id
  name       = "example"
}

data "stackit_object_storage_bucket" "example" {
  depends_on = [stackit_object_storage_bucket.example]
  project_id = var.project_id
  name       = "example"
}
