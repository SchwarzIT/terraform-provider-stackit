resource "stackit_project" "example" {
  name        = "example"
  billing_ref = var.project_billing_ref
  owner_id    = var.project_owner
}
