resource "stackit_mariadb_instance" "example" {
  name       = "example"
  project_id = "example"
  version    = "10.4"
  plan       = "stackit-mariadb-single-small"
}
