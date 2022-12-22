package credentialsgroup_test

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

func TestAcc_ObjectStorageCredentialsGroup(t *testing.T) {
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
			// check minimal configuration
			{
				Config: config(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("stackit_object_storage_credentials_group.example", "id"),
					resource.TestCheckResourceAttrSet("stackit_object_storage_credentials_group.example", "urn"),
					resource.TestCheckResourceAttr("data.stackit_object_storage_credentials_group.example", "name", name),
					resource.TestCheckResourceAttr("data.stackit_object_storage_credentials_group.example", "object_storage_project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckTypeSetElemAttrPair("stackit_object_storage_credentials_group.example", "id", "data.stackit_object_storage_credentials_group.example", "id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_object_storage_credentials_group.example", "urn", "data.stackit_object_storage_credentials_group.example", "urn"),
					resource.TestCheckTypeSetElemAttrPair("stackit_object_storage_credentials_group.example", "name", "data.stackit_object_storage_credentials_group.example", "name"),
					resource.TestCheckTypeSetElemAttrPair("stackit_object_storage_credentials_group.example", "object_storage_project_id", "data.stackit_object_storage_credentials_group.example", "object_storage_project_id"),
				),
			},
		},
	})
}

func config(name string) string {
	return fmt.Sprintf(`
	resource "stackit_object_storage_project" "example" {
		project_id         = "%s"
	}

	resource "stackit_object_storage_credentials_group" "example" {
		object_storage_project_id 	= stackit_object_storage_project.example.id
		name	   					= "%s"
	}

	data "stackit_object_storage_credentials_group" "example" {
		object_storage_project_id 	= stackit_object_storage_project.example.id
		id		  					= stackit_object_storage_credentials_group.example.id
	}
	`,
		common.GetAcceptanceTestsProjectID(),
		name,
	)
}
