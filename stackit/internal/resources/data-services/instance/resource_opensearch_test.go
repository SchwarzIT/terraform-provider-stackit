package instance_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const opensearch_inst_run_this_test = false

func TestAcc_ResourceOpensearchInstanceJob(t *testing.T) {
	if !common.ShouldAccTestRun(opensearch_inst_run_this_test) {
		t.Skip()
		return
	}

	name := "odjtest-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)
	plan1 := "stackit-opensearch-1.4.10-single"
	planID1 := "24615c29-99e8-4cc2-bcc3-ad7f45a5d46f"
	plan2 := "stackit-opensearch-2.4.10-single"
	planID2 := "f97a4935-0a77-4939-bfd1-33ba2e1f2b36"
	version := "2"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: configInstopensearch(name, plan1, version),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_opensearch_instance.example", "name", name),
					resource.TestCheckResourceAttr("stackit_opensearch_instance.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttr("stackit_opensearch_instance.example", "version", version),
					resource.TestCheckResourceAttr("stackit_opensearch_instance.example", "plan", plan1),
					resource.TestCheckResourceAttr("stackit_opensearch_instance.example", "plan_id", planID1),
					resource.TestCheckResourceAttrSet("stackit_opensearch_instance.example", "id"),
					resource.TestCheckResourceAttrSet("stackit_opensearch_instance.example", "dashboard_url"),
					resource.TestCheckResourceAttrSet("stackit_opensearch_instance.example", "cf_guid"),
					resource.TestCheckResourceAttrSet("stackit_opensearch_instance.example", "cf_space_guid"),
				),
			},
			// check update plan
			{
				Config: configInstopensearch(name, plan2, version),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_opensearch_instance.example", "name", name),
					resource.TestCheckResourceAttr("stackit_opensearch_instance.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttr("stackit_opensearch_instance.example", "version", version),
					resource.TestCheckResourceAttr("stackit_opensearch_instance.example", "plan", plan2),
					resource.TestCheckResourceAttr("stackit_opensearch_instance.example", "plan_id", planID2),
					resource.TestCheckResourceAttrSet("stackit_opensearch_instance.example", "id"),
					resource.TestCheckResourceAttrSet("stackit_opensearch_instance.example", "dashboard_url"),
					resource.TestCheckResourceAttrSet("stackit_opensearch_instance.example", "cf_guid"),
					resource.TestCheckResourceAttrSet("stackit_opensearch_instance.example", "cf_space_guid"),
				),
			},
			// test import
			{
				ResourceName: "stackit_opensearch_instance.example",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r, ok := s.RootModule().Resources["stackit_opensearch_instance.example"]
					if !ok {
						return "", errors.New("couldn't find resource stackit_opensearch_instance.example")
					}
					id, ok := r.Primary.Attributes["id"]
					if !ok {
						return "", errors.New("couldn't find attribute id")
					}

					return fmt.Sprintf("%s,%s", common.GetAcceptanceTestsProjectID(), id), nil
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func configInstopensearch(name, plan, version string) string {
	return fmt.Sprintf(`
	resource "stackit_opensearch_instance" "example" {
		name       = "%s"
		project_id = "%s"
		version    = "%s"
		plan       = "%s"
	  }
	  
	  `,
		name,
		common.GetAcceptanceTestsProjectID(),
		version,
		plan,
	)
}
