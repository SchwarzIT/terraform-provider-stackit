resource "stackit_opensearch_instance" "example" {
  name       = "example"
  project_id = "example"
}

resource "stackit_opensearch_credential" "example" {
  project_id  = "example"
  instance_id = stackit_opensearch_instance.example.id
}
