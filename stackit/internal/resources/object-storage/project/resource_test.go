package project_test

import (
	"fmt"
	"testing"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const run_this_test = false

func TestAcc_object_storageProject(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_object_storage_project.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttr("stackit_object_storage_project.example", "id", common.GetAcceptanceTestsProjectID()),
				),
			},
			// test import
			{
				ResourceName:            "stackit_object_storage_project.example",
				ImportStateId:           fmt.Sprintf("%s", common.GetAcceptanceTestsProjectID()),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func config() string {
	return fmt.Sprintf(`
resource "stackit_object_storage_project" "example" {
	project_id         = "%s"
}
	  `,
		common.GetAcceptanceTestsProjectID(),
	)
}
