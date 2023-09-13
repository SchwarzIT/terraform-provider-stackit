terraform {
  required_providers {
    stackit = {
      source = "github.com/schwarzit/stackit"
    }
    openstack = {
      source = "terraform-provider-openstack/openstack"
    }
  }
}

provider "stackit" {}

# Create a token for the OpenStack provider on your project's Infrastructure API
provider "openstack" {
  tenant_id        = "{OpenStack project ID}"
  tenant_name      = "{OpenStack project name}"
  user_name        = "{Token name}"
  user_domain_name = "portal_mvp"
  password         = "{Token password}"
  region           = "RegionOne"
  auth_url         = "https://keystone.api.iaas.eu01.stackit.cloud/v3"
}
