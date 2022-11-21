resource "stackit_logme_instance" "example" {
  name       = "example"
  project_id = "example"
  version    = "LogMe"
  plan       = "stackit-logme-single-small-non-ssl"
}
