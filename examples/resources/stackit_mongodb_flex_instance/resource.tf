resource "stackit_mongodb_flex_instance" "example" {
  name         = "example"
  project_id   = "example"
  machine_type = "1.1"
  acl = [
    "193.148.160.0/19",
    "45.129.40.0/21",
    "45.135.244.0/22"
  ]
}
