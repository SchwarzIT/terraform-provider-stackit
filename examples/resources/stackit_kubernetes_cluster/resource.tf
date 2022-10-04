resource "stackit_project" "example" {
  name        = "example"
  billing_ref = var.project_billing_ref
  owner       = var.project_owner
}


resource "stackit_kubernetes_cluster" "example" {
  name               = "example"
  project_id         = stackit_project.example.id
  kubernetes_version = "1.23.12"

  node_pool {
    name         = "example-np"
    machine_type = "c1.2"
    zones        = ["eu01-m"]
  }
}
