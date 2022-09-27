resource "stackit_project" "example" {
  name        = "example"
  billing_ref = var.project_billing_ref
  owner       = var.project_owner
}
