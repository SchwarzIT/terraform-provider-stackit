package instance_test

import (
	"fmt"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

const run_this_test = false

func TestAcc_kubernetes(t *testing.T) {
	if !run_this_test {
		t.Skip()
		return
	}
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: config(name, "instancepl"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.stackit_argus_instance", "name", name),
					resource.TestCheckResourceAttr("data.stackit_argus_instance", "project_id", common.ACC_TEST_PROJECT_ID),
					resource.TestCheckResourceAttr("data.stackit_argus_instance", "grafana.enable_public_access", "true"),
					resource.TestCheckResourceAttr("data.stackit_argus_instance", "metrics.retention_days", "60"),
					resource.TestCheckResourceAttr("data.stackit_argus_instance", "metrics.retention_days_5m_downsampling", "20"),
					resource.TestCheckResourceAttr("data.stackit_argus_instance", "metrics.retention_days_1h_downsampling", "10"),
				),
			},
		},
	})
}

func config(name, plan string) string {
	return fmt.Sprintf(`
resource "stackit_argus_instance" "example" {
	project_id = "%s"
	name       = "%s"
	plan       = "instancepl"
	grafana	   = {
		enable_public_access = true
	}
	metrics	   = {
		retention_days 				   = 60
		retention_days_5m_downsampling = 20
		retention_days_1h_downsampling = 10
	}
}
	data "data.stackit_argus_instance" "example" {
	depends_on = [stackit_argus_instance.example]
	project_id = "%s"
	name       = "%s"
}
	  `,
		common.ACC_TEST_PROJECT_ID,
		name,
		common.ACC_TEST_PROJECT_ID,
		name,
	)
}
