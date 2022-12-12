package credential_test

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

func TestAcc_ObjectStorageCredential(t *testing.T) {
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
					resource.TestCheckTypeSetElemAttrPair("stackit_object_storage_credential.example", "display_name", "data.stackit_object_storage_credential.ex1", "display_name"),
					resource.TestCheckTypeSetElemAttrPair("stackit_object_storage_credential.example", "id", "data.stackit_object_storage_credential.ex2", "id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_object_storage_credential.example", "expiry", "data.stackit_object_storage_credential.ex1", "expiry"),
					resource.TestCheckTypeSetElemAttrPair("stackit_object_storage_credential.example", "expiry", "data.stackit_object_storage_credential.ex2", "expiry"),
					resource.TestCheckTypeSetElemAttrPair("stackit_object_storage_credential.example", "project_id", "data.stackit_object_storage_credential.ex1", "project_id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_object_storage_credential.example", "project_id", "data.stackit_object_storage_credential.ex2", "project_id"),
				),
			},
		},
	})
}

func config() string {
	return fmt.Sprintf(`
	resource "stackit_object_storage_credential" "example" {
		project_id = "%s"
	}


	data "stackit_object_storage_credential" "ex1" {
		project_id = "%s"
		id		   = stackit_object_storage_credential.example.id
	}

	data "stackit_object_storage_credential" "ex2" {
		project_id   = "%s"
		display_name = stackit_object_storage_credential.example.display_name
	}

	`,
		common.GetAcceptanceTestsProjectID(),
		common.GetAcceptanceTestsProjectID(),
		common.GetAcceptanceTestsProjectID(),
	)
}
