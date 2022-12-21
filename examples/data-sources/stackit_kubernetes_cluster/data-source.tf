resource "stackit_project" "example" {
  name        = "example"
  billing_ref = var.project_billing_ref
  owner       = var.project_owner
}

data "stackit_kubernetes_project" "example" {
  project_id = stackit_project.example.id
}

data "stackit_kubernetes_cluster" "example" {
  name                  = "example"
  kubernetes_project_id = data.stackit_kubernetes_project.example.id
}
