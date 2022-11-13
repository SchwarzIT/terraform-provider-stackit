resource "stackit_elasticsearch_instance" "example" {
  name       = "some_name_2"
  project_id = var.project_id
  version    = "7"
  plan       = "stackit-elasticsearch-single-medium"
}
