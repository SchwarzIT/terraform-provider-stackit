resource "stackit_kubernetes_cluster" "example" {
  name               = "example"
  project_id         = stackit_project.example.id
  kubernetes_version = "1.23"

  node_pools = [{
    name         = "example"
    machine_type = "c1.2"
  }]
}
