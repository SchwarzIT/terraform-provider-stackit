resource "stackit_kubernetes_project" "example" {
  project_id = "example"
}

resource "stackit_kubernetes_cluster" "example" {
  name                  = "example"
  kubernetes_project_id = stackit_kubernetes_project.example.id

  node_pools = [{
    name         = "example"
    machine_type = "c1.2"
  }]
}
