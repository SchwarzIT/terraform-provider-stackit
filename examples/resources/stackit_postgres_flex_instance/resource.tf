resource "stackit_postgres_flex_instance" "example" {
  name         = "example"
  project_id   = "example"
  machine_type = "c1.2"
  version      = "14"
  replicas     = 1
}
