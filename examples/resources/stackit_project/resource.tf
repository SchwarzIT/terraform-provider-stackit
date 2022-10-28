resource "stackit_project" "example" {
  name                = "example"
  parent_container_id = "parent-contaier-id"
  billing_ref         = var.project_billing_ref
  owner_email         = var.project_owner_email
}
