resource "stackit_redis_instance" "example" {
  name       = "example"
  project_id = "example"
  version    = "3.7"
  plan       = "stackit-redis-single-small"
}
