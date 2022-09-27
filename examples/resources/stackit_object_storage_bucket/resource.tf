resource "stackit_project" "example" {
  name        = "example"
  billing_ref = var.project_billing_ref
  owner       = var.project_owner
}

resource "stackit_object_storage_bucket" "example" {
  name       = "example"
  project_id = stackit_project.example.id
}
