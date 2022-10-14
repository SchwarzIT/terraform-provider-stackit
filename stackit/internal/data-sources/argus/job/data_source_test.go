package job_test

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

func TestAcc_ArgusJob(t *testing.T) {
	if !run_this_test {
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
				Config: config(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.stackit_argus_job.example", "name", name),
					resource.TestCheckResourceAttr("data.stackit_argus_job.example", "project_id", common.ACC_TEST_PROJECT_ID),

					resource.TestCheckResourceAttr("data.stackit_argus_job.example", "metrics_path", "/metrics"),
					resource.TestCheckResourceAttr("data.stackit_argus_job.example", "scheme", "https"),
					resource.TestCheckResourceAttr("data.stackit_argus_job.example", "scrape_interval", "5m"),
					resource.TestCheckResourceAttr("data.stackit_argus_job.example", "scrape_timeout", "2m"),
					resource.TestCheckResourceAttr("data.stackit_argus_job.example", "saml2.enable_url_parameters", "true"),
					resource.TestCheckResourceAttr("data.stackit_argus_job.example", "targets.0.urls.0", "url1"),
					resource.TestCheckResourceAttr("data.stackit_argus_job.example", "targets.0.urls.1", "url2"),
				),
			},
		},
	})
}

func config(name string) string {
	return fmt.Sprintf(`
resource "stackit_argus_instance" "example" {
	project_id = "%s"
	name       = "example" 
	plan       = "Monitoring-Medium-EU01"
}

resource "stackit_argus_job" "example" {
	name              = "%s"
	project_id 		  = "%s"
	argus_instance_id = stackit_argus_instance.example.id
	scrape_interval   = "5m"
	scrape_timeout 	  = "2m"
	targets = [
	  {
		urls = ["url1", "url2"]
	  }
	]
}

data "stackit_argus_job" "example" {	
	depends_on				   = [stackit_argus_job.example]
	project_id 				   = "%s"
	name                       = stackit_argus_job.example.name
	argus_instance_id		   = stackit_argus_instance.example.id
}

	  `,
		common.ACC_TEST_PROJECT_ID,
		name,
		common.ACC_TEST_PROJECT_ID,
		common.ACC_TEST_PROJECT_ID,
	)
}
