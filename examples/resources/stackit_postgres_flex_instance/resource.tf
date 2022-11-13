resource "stackit_postgres_flex_instance" "example" {
  name       = "example"
  project_id = var.project_id
  version    = "7"
  plan       = "stackit-elasticsearch-single-small"
}
