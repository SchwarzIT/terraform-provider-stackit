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

const run_this_test = false

func TestAcc_ElasticSearchInstance(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}

	name := "odjtest-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)
	plan := "stackit-elasticsearch-single-small"
	planID := "a59cf7bb-ae64-4f63-8503-fca8c936bf0c"
	version := "7"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(name, version, plan),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.stackit_elasticsearch_instance.example", "name", name),
					resource.TestCheckResourceAttr("data.stackit_elasticsearch_instance.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttr("data.stackit_elasticsearch_instance.example", "version", version),
					resource.TestCheckResourceAttr("data.stackit_elasticsearch_instance.example", "plan", plan),
					resource.TestCheckResourceAttr("data.stackit_elasticsearch_instance.example", "plan_id", planID),
					resource.TestCheckResourceAttrSet("data.stackit_elasticsearch_instance.example", "id"),
					resource.TestCheckResourceAttrSet("data.stackit_elasticsearch_instance.example", "dashboard_url"),
					resource.TestCheckResourceAttrSet("data.stackit_elasticsearch_instance.example", "cf_guid"),
					resource.TestCheckResourceAttrSet("data.stackit_elasticsearch_instance.example", "cf_space_guid"),
					resource.TestCheckResourceAttrSet("data.stackit_elasticsearch_instance.example", "cf_organization_guid"),
					resource.TestCheckTypeSetElemAttrPair("stackit_elasticsearch_instance.example", "id", "data.stackit_elasticsearch_instance.example", "id"),
				),
			},
		},
	})
}

func config(name, version, plan string) string {
	return fmt.Sprintf(`
	resource "stackit_elasticsearch_instance" "example" {
		name       = "%s"
		project_id = "%s"
		version    = "%s"
		plan       = "%s"
	  }


	
	  data "stackit_elasticsearch_instance" "example" {
		depends_on = [stackit_elasticsearch_instance.example]
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
