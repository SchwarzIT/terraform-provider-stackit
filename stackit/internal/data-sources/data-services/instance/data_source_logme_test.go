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

const logme_inst_run_this_test = false

func TestAcc_DataSourceLogMeInstanceJob(t *testing.T) {
	if !common.ShouldAccTestRun(logme_inst_run_this_test) {
		t.Skip()
		return
	}

	name := "odjtest-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)
	plan := "stackit-logme2-1.4.10-single"
	planID := "7a54492c-8a2e-4d3c-b6c2-a4f20cb65912"
	version := "2"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: configLogMe(name, version, plan),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.stackit_logme_instance.example", "name", name),
					resource.TestCheckResourceAttr("data.stackit_logme_instance.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttr("data.stackit_logme_instance.example", "version", version),
					resource.TestCheckResourceAttr("data.stackit_logme_instance.example", "plan", plan),
					resource.TestCheckResourceAttr("data.stackit_logme_instance.example", "plan_id", planID),
					resource.TestCheckResourceAttrSet("data.stackit_logme_instance.example", "id"),
					resource.TestCheckResourceAttrSet("data.stackit_logme_instance.example", "dashboard_url"),
					resource.TestCheckResourceAttrSet("data.stackit_logme_instance.example", "cf_guid"),
					resource.TestCheckResourceAttrSet("data.stackit_logme_instance.example", "cf_space_guid"),
					resource.TestCheckTypeSetElemAttrPair("stackit_logme_instance.example", "id", "data.stackit_logme_instance.example", "id"),
				),
			},
		},
	})
}

func configLogMe(name, version, plan string) string {
	return fmt.Sprintf(`
	resource "stackit_logme_instance" "example" {
		name       = "%s"
		project_id = "%s"
		version    = "%s"
		plan       = "%s"
	  }

	
	  data "stackit_logme_instance" "example" {
		depends_on = [stackit_logme_instance.example]
		name       = "%s"
		project_id = "%s"
	  }

	`,
		name,
		common.GetAcceptanceTestsProjectID(),
		version,
		plan,
		name,
		common.GetAcceptanceTestsProjectID(),
	)
}
