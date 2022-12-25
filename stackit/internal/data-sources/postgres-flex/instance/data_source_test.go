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

func TestAcc_PostgresFlexInstance(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}

	name := "odjtest-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)
	version := "14"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(name, version),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.stackit_postgres_flex_instance.example", "name", name),
					resource.TestCheckResourceAttr("data.stackit_postgres_flex_instance.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttr("data.stackit_postgres_flex_instance.example", "version", version),
					resource.TestCheckResourceAttr("stackit_postgres_flex_instance.example", "machine_type", "c1.2"),
					resource.TestCheckResourceAttr("stackit_postgres_flex_instance.example", "replicas", "1"),
					resource.TestCheckResourceAttr("stackit_postgres_flex_instance.example", "storage.class", "premium-perf6-stackit"),
					resource.TestCheckResourceAttr("stackit_postgres_flex_instance.example", "storage.size", "30"),
					resource.TestCheckResourceAttrSet("stackit_postgres_flex_instance.example", "id"),
				),
			},
		},
	})
}

func config(name, version string) string {
	return fmt.Sprintf(`
	resource "stackit_postgres_flex_instance" "example" {
		name         = "%s"
		project_id   = "%s"
		machine_type = "c1.2"
		version      = "%s"
		replicas     = 1
		storage 	 = {
			class = "premium-perf6-stackit"
			size  = 30
		}
	}  

	
	  data "stackit_postgres_flex_instance" "example" {
		depends_on = [stackit_postgres_flex_instance.example]
		name       = "%s"
		project_id = "%s"
	  }

	`,
		name,
		common.GetAcceptanceTestsProjectID(),
		version,
		name,
		common.GetAcceptanceTestsProjectID(),
	)
}
