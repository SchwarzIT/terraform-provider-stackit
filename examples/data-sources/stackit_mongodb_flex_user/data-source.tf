resource "stackit_mongodb_flex_instance" "example" {
  project_id   = var.project_id
  name         = "example"
  machine_type = "1.1"
}
resource "stackit_mongodb_flex_user" "example" {
  project_id  = var.project_id
  instance_id = stackit_mongodb_flex_instance.example.id
}

data "stackit_mongodb_flex_user" "example" {
  project_id  = var.project_id
  instance_id = stackit_mongodb_flex_instance.example.id
  id          = stackit_mongodb_flex_user.example.id
}
