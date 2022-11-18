resource "stackit_redis_instance" "example" {
  name       = "example"
  project_id = "example"
  version    = "6.0"
  plan       = "stackit-redis-single-small"
}
