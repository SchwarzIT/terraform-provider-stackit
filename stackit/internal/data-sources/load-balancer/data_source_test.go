package loadbalancer_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type openstack struct {
	tenantID   string
	tenantName string
	userName   string
	password   string
}

const run_this_test = false

func TestAcc_LoadBalancer(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}
	os := openstack{
		tenantID:   os.Getenv("OS_TENANT_ID"),
		tenantName: os.Getenv("OS_TENANT_NAME"),
		userName:   os.Getenv("OS_USERNAME"),
		password:   os.Getenv("OS_PASSWORD"),
	}
	projectID := "8a2d2862-ac85-4084-8144-4c72d92ddcdd"
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		ExternalProviders: map[string]resource.ExternalProvider{
			"openstack": {
				VersionConstraint: "= 1.52.1",
				Source:            "terraform-provider-openstack/openstack",
			},
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(name, projectID, os),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("stackit_load_balancer.example", "project_id", "data.stackit_load_balancer.example", "project_id"),
					resource.TestCheckResourceAttrPair("stackit_load_balancer.example", "id", "data.stackit_load_balancer.example", "id"),
					resource.TestCheckResourceAttrPair("stackit_load_balancer.example", "external_address", "data.stackit_load_balancer.example", "external_address"),
					resource.TestCheckResourceAttrPair("stackit_load_balancer.example", "private_network_only", "data.stackit_load_balancer.example", "private_network_only"),
					resource.TestCheckResourceAttrPair("stackit_load_balancer.example", "name", "data.stackit_load_balancer.example", "name"),
					resource.TestCheckResourceAttrPair("stackit_load_balancer.example", "target_pools.0.name", "data.stackit_load_balancer.example", "target_pools.0.name"),
					resource.TestCheckResourceAttrPair("stackit_load_balancer.example", "target_pools.0.target_port", "data.stackit_load_balancer.example", "target_pools.0.target_port"),
					resource.TestCheckResourceAttrPair("stackit_load_balancer.example", "target_pools.0.targets.0.display_name", "data.stackit_load_balancer.example", "target_pools.0.targets.0.display_name"),
					resource.TestCheckResourceAttrPair("stackit_load_balancer.example", "target_pools.0.targets.0.ip_address", "data.stackit_load_balancer.example", "target_pools.0.targets.0.ip_address"),
					resource.TestCheckResourceAttrPair("stackit_load_balancer.example", "listeners.0.display_name", "data.stackit_load_balancer.example", "listeners.0.display_name"),
					resource.TestCheckResourceAttrPair("stackit_load_balancer.example", "listeners.0.port", "data.stackit_load_balancer.example", "listeners.0.port"),
					resource.TestCheckResourceAttrPair("stackit_load_balancer.example", "listeners.0.protocol", "data.stackit_load_balancer.example", "listeners.0.protocol"),
					resource.TestCheckResourceAttrPair("stackit_load_balancer.example", "listeners.0.target_pool", "data.stackit_load_balancer.example", "listeners.0.target_pool"),
					resource.TestCheckResourceAttrPair("stackit_load_balancer.example", "networks.0.network_id", "data.stackit_load_balancer.example", "networks.0.network_id"),
				),
			},
		},
	})
}

func config(name, projectID string, os openstack) string {
	return fmt.Sprintf(`
	resource "stackit_load_balancer" "example" {
		project_id           = "%s"
		name                 = "%s"
		external_address     = openstack_networking_floatingip_v2.example_ip.address
		private_network_only = false
		target_pools = [{
			name        = "example-target-pool"
			target_port = 80
			targets = [{
				display_name = "example-target"
				ip_address   = openstack_compute_instance_v2.example.network.0.fixed_ip_v4
			}]
		}]
		listeners = [{
			display_name = "example-listener"
			port         = 80
			protocol     = "PROTOCOL_TCP"
			target_pool  = "example-target-pool"
		}]
		networks = [
			{ network_id = openstack_networking_network_v2.example.id }
		]
	}

	data "stackit_load_balancer" "example" {
		depends_on = [stackit_load_balancer.example]
		project_id = stackit_load_balancer.example.project_id
		name       = stackit_load_balancer.example.name
	}

%s
  
	  `, projectID, name, supportingInfra(name, os))
}

func supportingInfra(name string, os openstack) string {
	return fmt.Sprintf(`

	provider "openstack" {
		tenant_id        = "%s"
		tenant_name      = "%s"
		user_name        = "%s"
		user_domain_name = "portal_mvp"
		password         = "%s"
		region           = "RegionOne"
		auth_url         = "https://keystone.api.iaas.eu01.stackit.cloud/v3"
	}

	# Create a network
	resource "openstack_networking_network_v2" "example" {
	  name = "%s_network"
	}
	
	resource "openstack_networking_subnet_v2" "example" {
	  name            = "%s_subnet"
	  cidr            = "192.168.0.0/25"
	  dns_nameservers = ["8.8.8.8"]
	  network_id      = openstack_networking_network_v2.example.id
	}
	
	data "openstack_networking_network_v2" "public" {
	  name = "floating-net"
	}
	
	resource "openstack_networking_floatingip_v2" "example_ip" {
	  pool = data.openstack_networking_network_v2.public.name
	}
	
	# Create an instance
	data "openstack_compute_flavor_v2" "example" {
	  name = "g1.1"
	}


	# Create instance
	resource "openstack_compute_instance_v2" "example" {
		depends_on      = [openstack_networking_subnet_v2.example]
		name            = "%s_instance"
		flavor_id       = data.openstack_compute_flavor_v2.example.id
		admin_pass      = "example"
		security_groups = ["default"]

		block_device {
			uuid                  = "4364cdb2-dacd-429b-803e-f0f7cfde1c24" // Ubuntu 22.04
			volume_size           = 32
			source_type           = "image"
			destination_type      = "volume"
			delete_on_termination = true
		}

		network {
			name = openstack_networking_network_v2.example.name
		}

		lifecycle {
			ignore_changes = [security_groups]
		}
	}
	
	resource "openstack_networking_router_v2" "example_router" {
	  name                = "%s_router"
	  admin_state_up      = "true"
	  external_network_id = data.openstack_networking_network_v2.public.id
	}
	
	resource "openstack_networking_router_interface_v2" "example_interface" {
	  router_id = openstack_networking_router_v2.example_router.id
	  subnet_id = openstack_networking_subnet_v2.example.id
	}
	
`, os.tenantID, os.tenantName, os.userName, os.password, name, name, name, name)
}
