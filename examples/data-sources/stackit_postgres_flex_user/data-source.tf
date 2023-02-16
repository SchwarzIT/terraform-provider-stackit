resource "stackit_postgres_flex_instance" "example" {
  name         = "example"
  project_id   = var.project_id
  machine_type = "c1.2"
  version      = "14"
}

resource "stackit_postgres_flex_user" "example" {
  project_id  = var.project_id
  instance_id = stackit_postgres_flex_instance.example.id
}

data "stackit_postgres_flex_user" "example" {
  project_id  = var.project_id
  instance_id = stackit_postgres_flex_instance.example.id
  id          = stackit_postgres_flex_user.example.id
}
