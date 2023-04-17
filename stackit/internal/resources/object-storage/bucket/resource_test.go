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

const run_this_test = true

func TestAcc_ObjectStorageBucket(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}

	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	newName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_object_storage_bucket.example", "name", name),
					resource.TestCheckResourceAttr("stackit_object_storage_bucket.example", "object_storage_project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttrSet("stackit_object_storage_bucket.example", "region"),
					resource.TestCheckResourceAttrSet("stackit_object_storage_bucket.example", "host_style_url"),
					resource.TestCheckResourceAttrSet("stackit_object_storage_bucket.example", "path_style_url"),
				),
			},
			// check minimal configuration
			{
				Config: config(newName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_object_storage_bucket.example", "name", newName),
					resource.TestCheckResourceAttr("stackit_object_storage_bucket.example", "object_storage_project_id", common.GetAcceptanceTestsProjectID()),
				),
			},
			// test import
			{
				ResourceName:      "stackit_object_storage_bucket.example",
				ImportStateId:     fmt.Sprintf("%s,%s", common.GetAcceptanceTestsProjectID(), newName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func config(name string) string {
	return fmt.Sprintf(`
resource "stackit_object_storage_project" "example" {
	project_id         = "%s"
}

resource "stackit_object_storage_bucket" "example" {
	object_storage_project_id = stackit_object_storage_project.example.id
    name      				  = "%s"
}
	  `,
		common.GetAcceptanceTestsProjectID(),
		name,
	)
}
