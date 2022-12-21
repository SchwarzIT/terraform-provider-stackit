resource "stackit_kubernetes_project" "example" {
  project_id = "example"
}

data "stackit_kubernetes_project" "example" {
  depends_on = [stackit_kubernetes_project.example]
  project_id = "example"
}
