resource "stackit_redis_instance" "example" {
  name       = "example"
  project_id = "example"
}
resource "stackit_redis_credential" "example" {
  project_id  = "example"
  instance_id = stackit_redis_instance.example.id
}
