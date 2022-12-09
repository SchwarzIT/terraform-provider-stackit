resource "stackit_mariadb_instance" "example" {
  name       = "example"
  project_id = "example"
}

resource "stackit_mariadb_credential" "example" {
  project_id  = "example"
  instance_id = stackit_mariadb_instance.example.id
}
