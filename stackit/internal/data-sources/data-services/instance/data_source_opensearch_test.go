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

const opensearch_inst_run_this_test = true

func TestAcc_DataSourceOpensearchInstanceJob(t *testing.T) {
	if !common.ShouldAccTestRun(opensearch_inst_run_this_test) {
		t.Skip()
		return
	}

	name := "odjtest-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: configopensearch(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.stackit_opensearch_instance.example", "name", name),
					resource.TestCheckResourceAttr("data.stackit_opensearch_instance.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckTypeSetElemAttrPair("stackit_opensearch_instance.example", "id", "data.stackit_opensearch_instance.example", "id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_opensearch_instance.example", "name", "data.stackit_opensearch_instance.example", "name"),
					resource.TestCheckTypeSetElemAttrPair("stackit_opensearch_instance.example", "plan", "data.stackit_opensearch_instance.example", "plan"),
					resource.TestCheckTypeSetElemAttrPair("stackit_opensearch_instance.example", "plan_id", "data.stackit_opensearch_instance.example", "plan_id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_opensearch_instance.example", "version", "data.stackit_opensearch_instance.example", "version"),
					resource.TestCheckTypeSetElemAttrPair("stackit_opensearch_instance.example", "dashboard_url", "data.stackit_opensearch_instance.example", "dashboard_url"),
					resource.TestCheckTypeSetElemAttrPair("stackit_opensearch_instance.example", "cf_guid", "data.stackit_opensearch_instance.example", "cf_guid"),
					resource.TestCheckTypeSetElemAttrPair("stackit_opensearch_instance.example", "cf_space_guid", "data.stackit_opensearch_instance.example", "cf_space_guid"),
				),
			},
		},
	})
}

func configopensearch(name string) string {
	return fmt.Sprintf(`
	resource "stackit_opensearch_instance" "example" {
		name       = "%s"
		project_id = "%s"
	  }

	
	  data "stackit_opensearch_instance" "example" {
		depends_on = [stackit_opensearch_instance.example]
		name       = "%s"
		project_id = "%s"
	  }

	`,
		name,
		common.GetAcceptanceTestsProjectID(),
		name,
		common.GetAcceptanceTestsProjectID(),
	)
}
