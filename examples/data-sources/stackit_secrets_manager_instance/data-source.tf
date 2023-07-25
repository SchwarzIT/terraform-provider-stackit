resource "stackit_secrets_manager_instance" "example" {
  project_id = var.project_id
  name       = "example"
}

data "stackit_secrets_manager_instance" "example" {
  project_id = var.project_id
  id         = stackit_secrets_manager_instance.example.id
}
