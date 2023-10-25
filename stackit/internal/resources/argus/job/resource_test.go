package job_test

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

const run_this_test = false

func TestAcc_ArgusJob(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}

	name := "e1" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_argus_job.example", "name", "example"),
					resource.TestCheckResourceAttr("stackit_argus_job.example", "project_id", common.GetAcceptanceTestsProjectID()),
				),
			},
			// check update
			{
				Config: config2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_argus_job.example", "name", "example"),
					resource.TestCheckResourceAttr("stackit_argus_job.example", "project_id", common.GetAcceptanceTestsProjectID()),
				),
			},
			// test import
			{
				ResourceName: "stackit_argus_job.example",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r, ok := s.RootModule().Resources["stackit_argus_instance.example"]
					if !ok {
						return "", errors.New("couldn't find resource stackit_argus_instance.example")
					}
					id, ok := r.Primary.Attributes["id"]
					if !ok {
						return "", errors.New("couldn't find attribute id")
					}

					return fmt.Sprintf("%s,%s,%s", common.GetAcceptanceTestsProjectID(), id, "example"), nil
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func config(name string) string {
	return fmt.Sprintf(`
resource "stackit_argus_instance" "example" {
	project_id = "%s"
	name       = "%s"
	plan       = "Monitoring-Medium-EU01"
}

resource "stackit_argus_job" "example" {
	name              = "example"
	project_id 		  = "%s"
	argus_instance_id = stackit_argus_instance.example.id
	targets = [
	  {
		urls = ["url1", "url2"]
	  }
	]

	saml2 = {
	  enable_url_parameters = true
	}
}
	  `,
		common.GetAcceptanceTestsProjectID(),
		name,
		common.GetAcceptanceTestsProjectID(),
	)
}

func config2(name string) string {
	return fmt.Sprintf(`
resource "stackit_argus_instance" "example" {
	project_id = "%s"
	name       = "%s"
	plan       = "Monitoring-Medium-EU01"
}

resource "stackit_argus_job" "example" {
	name              = "example"
	project_id 		  = "%s"
	argus_instance_id = stackit_argus_instance.example.id
	targets = [
	  {
		urls = ["url3", "url4"]
	  }
	]
}
	
	  `,
		common.GetAcceptanceTestsProjectID(),
		name,
		common.GetAcceptanceTestsProjectID(),
	)
}
