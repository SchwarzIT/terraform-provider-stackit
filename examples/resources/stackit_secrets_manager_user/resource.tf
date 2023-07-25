resource "stackit_secrets_manager_instance" "example" {
  project_id = var.project_id
  name       = "example"
}

resource "stackit_secrets_manager_user" "example" {
  project_id    = var.project_id
  instance_id   = stackit_secrets_manager_instance.example.id
  description   = "example"
  write_enabled = true
}
