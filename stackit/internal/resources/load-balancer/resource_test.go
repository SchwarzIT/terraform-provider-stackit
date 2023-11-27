package loadbalancer_test

import (
	"fmt"
	"testing"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const run_this_test = true

func TestAcc_LoadBalancer(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}
	projectID := "8a2d2862-ac85-4084-8144-4c72d92ddcdd"
	name := acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(projectID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_load_balancer.example", "project_id", projectID),
					resource.TestCheckResourceAttr("stackit_load_balancer.example", "id", fmt.Sprintf("example-%s", name)),
				),
			},
			// test import
			{
				ResourceName:      "stackit_load_balancer.example",
				ImportStateId:     fmt.Sprintf("%s,%s", projectID, fmt.Sprintf("example-%s", name)),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func config(projectID, name string) string {
	return fmt.Sprintf(`
	resource "stackit_load_balancer" "example" {
		project_id           = "%s"
		name                 = "example-%s"
		external_address     = "45.129.45.188"
		private_network_only = false
		target_pools = [{
		  name        = "example-target-pool"
		  target_port = 80
		  targets = [{
			display_name = "example-target"
			ip_address   = "192.168.0.94"
		  }]
		}]
		listeners = [{
		  display_name = "example-listener"
		  port         = 80
		  protocol     = "PROTOCOL_TCP"
		  target_pool  = "example-target-pool"
		}]
		networks = [
		  { network_id = "ccff5c2e-4133-4dd2-aa3c-c3f5786e47be" }
		]
	  }
  
	  `, projectID, name)
}
