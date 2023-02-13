resource "stackit_mongodb_flex_instance" "example" {
  name         = "example"
  project_id   = var.project_id
  machine_type = "1.1"
}

resource "stackit_mongodb_flex_user" "example" {
  project_id  = var.project_id
  instance_id = stackit_mongodb_flex_instance.example.id
}

output "mongodb_username" {
  value = stackit_mongodb_flex_user.example.username
}

output "mongodb_password" {
  value     = stackit_mongodb_flex_user.example.password
  sensitive = true
}
