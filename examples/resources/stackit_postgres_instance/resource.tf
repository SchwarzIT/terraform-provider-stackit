resource "stackit_postgres_instance" "example" {
  name       = "example"
  project_id = "example"
  version    = "13"
  plan       = "stackit-postgresql-single-small"
}
