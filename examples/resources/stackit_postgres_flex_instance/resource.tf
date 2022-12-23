resource "stackit_postgres_flex_instance" "example" {
  name         = "example"
  project_id   = "example"
  machine_type = "c1.2"
}
