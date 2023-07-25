resource "stackit_secrets_manager_instance" "example" {
  project_id = var.project_id
  name       = "example"
}

resource "stackit_secrets_manager_user" "example" {
  project_id    = var.project_id
  instance_id   = stackit_secrets_manager_instance.example.id
  description   = "test"
  write_enabled = true
}

data "stackit_secrets_manager_user" "example" {
  project_id  = var.project_id
  instance_id = stackit_secrets_manager_instance.example.id
  username    = stackit_secrets_manager_user.example.username
}
