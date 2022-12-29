package instance_test

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

func TestAcc_ArgusInstances(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: config(name, "Monitoring-Medium-EU01"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.stackit_argus_instance.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttr("data.stackit_argus_instance.example", "name", name),
					resource.TestCheckResourceAttr("data.stackit_argus_instance.example", "grafana.enable_public_access", "true"),
					resource.TestCheckResourceAttrSet("data.stackit_argus_instance.example", "id"),
					resource.TestCheckResourceAttrSet("data.stackit_argus_instance.example", "metrics.retention_days"),
					resource.TestCheckResourceAttrSet("data.stackit_argus_instance.example", "metrics.retention_days_5m_downsampling"),
					resource.TestCheckResourceAttrSet("data.stackit_argus_instance.example", "metrics.retention_days_1h_downsampling"),
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
	plan       = "%s"
	grafana	   = {
		enable_public_access = true
	}
	metrics	   = {
		retention_days				   = 60
		retention_days_5m_downsampling = 20
		retention_days_1h_downsampling = 10
	}
}

data "stackit_argus_instance" "example" {
	project_id = "%s"
	id		   = stackit_argus_instance.example.id
}

	  `,
		common.GetAcceptanceTestsProjectID(),
		name,
		plan,
		common.GetAcceptanceTestsProjectID(),
	)
}
