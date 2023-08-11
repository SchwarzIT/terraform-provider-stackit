resource "stackit_kubernetes_cluster" "example" {
  name       = "example"
  project_id = var.project_id

  node_pools = [{
    name         = "example"
    machine_type = "c1.2"
  }]
}
