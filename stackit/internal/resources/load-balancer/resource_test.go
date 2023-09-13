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

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_load_balancer.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttrSet("stackit_load_balancer.example", "id"),
				),
			},
		},
	})
}

func config() string {
	return fmt.Sprintf(`
resource "stackit_load_balancer" "example" {
	project_id = "%s"
	name	   = "test"
}
	  `,
		common.GetAcceptanceTestsProjectID(),
	)
}
