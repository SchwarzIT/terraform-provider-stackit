package bucket_test

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

func TestAcc_ObjectStorageBucket(t *testing.T) {
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
			// check configuration
			{
				Config: config(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.stackit_object_storage_bucket.example", "name", name),
					resource.TestCheckResourceAttr("data.stackit_object_storage_bucket.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttrSet("data.stackit_object_storage_bucket.example", "region"),
					resource.TestCheckResourceAttrSet("data.stackit_object_storage_bucket.example", "host_style_url"),
					resource.TestCheckResourceAttrSet("data.stackit_object_storage_bucket.example", "path_style_url"),
					resource.TestCheckResourceAttrSet("data.stackit_object_storage_bucket.example", "id"),
				),
			},
		},
	})
}

func config(name string) string {
	return fmt.Sprintf(`
resource "stackit_object_storage_bucket" "example" {
	project_id = "%s"
    name       = "%s"
}

data "stackit_object_storage_bucket" "example" {
	depends_on = [stackit_object_storage_bucket.example]
	project_id = "%s"
    name       = "%s"
}
	  `,
		common.GetAcceptanceTestsProjectID(),
		name,
		common.GetAcceptanceTestsProjectID(),
		name,
	)
}
