resource "stackit_argus_instance" "example" {
  name       = "example"
  project_id = stackit_project.example.id
  plan       = "Monitoring-Medium-EU01"
}

resource "stackit_argus_job" "example" {
  name              = "example"
  project_id        = stackit_project.example.id
  argus_instance_id = stackit_argus_instance.argus.id
  targets = [
    {
      urls = ["url1", "url2"]
    }
  ]
}
