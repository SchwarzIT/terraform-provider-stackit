resource "stackit_opensearch_instance" "example" {
  name       = "example"
  project_id = "example"
}

data "stackit_opensearch_instance" "example" {
  name       = "example"
  project_id = "example"
}
