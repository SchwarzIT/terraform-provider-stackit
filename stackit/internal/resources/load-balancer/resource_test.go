package loadbalancer_test

import (
	"fmt"
	"testing"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const run_this_test = true

func TestAcc_LoadBalancer(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}
	projectID := "8a2d2862-ac85-4084-8144-4c72d92ddcdd"
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_load_balancer.example", "project_id", projectID),
					resource.TestCheckResourceAttr("stackit_load_balancer.example", "id", "example"),
				),
			},
			// test import
			{
				ResourceName:      "stackit_load_balancer.example",
				ImportStateId:     fmt.Sprintf("%s,%s", projectID, "example"),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func config() string {
	return `
resource "stackit_load_balancer" "example" {
	project_id = "8a2d2862-ac85-4084-8144-4c72d92ddcdd"
	name       = "example"
	target_pools = [{
	  name        = "example-target-pool"
	  target_port = 80
	  targets = [{
		display_name = "example-target"
		ip_address   = "192.168.0.112"
	  }]
	}]
	listeners = [{
	  display_name = "example-listener"
	  port         = 80
	  protocol     = "PROTOCOL_TCP"
	  target_pool  = "example-target-pool"
	}]
	networks = [
	  { network_id = "ab320bc4-71ea-4eed-aa74-ace03e1af597" }
	]
	external_address     = "193.148.170.115"
	private_network_only = false
  }
  
	  `
}
