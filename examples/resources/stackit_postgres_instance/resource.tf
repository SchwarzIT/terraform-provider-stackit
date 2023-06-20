resource "stackit_postgres_instance" "example" {
  name       = "example"
  project_id = var.project_id
}
