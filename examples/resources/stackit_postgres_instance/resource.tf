resource "stackit_postgres_instance" "example" {
  name       = "example"
  project_id = "example"
  version    = "11"
  plan       = "stackit-postgresql-single-small"
}
