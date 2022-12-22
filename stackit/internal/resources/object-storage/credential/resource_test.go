package credential_test

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

func TestAcc_ObjectStorageCredential(t *testing.T) {
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
				Config: config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("stackit_object_storage_credential.example", "id"),
					resource.TestCheckResourceAttrSet("stackit_object_storage_credential.example", "display_name"),
					resource.TestCheckResourceAttrSet("stackit_object_storage_credential.example", "access_key"),
					resource.TestCheckResourceAttrSet("stackit_object_storage_credential.example", "secret_access_key"),
				),
			},
			{
				Config: configWithGroup(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("stackit_object_storage_credential.example", "id"),
					resource.TestCheckResourceAttrSet("stackit_object_storage_credential.example", "display_name"),
					resource.TestCheckResourceAttrSet("stackit_object_storage_credential.example", "access_key"),
					resource.TestCheckResourceAttrSet("stackit_object_storage_credential.example", "secret_access_key"),
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
	`,
		common.GetAcceptanceTestsProjectID(),
	)
}

func configWithGroup(groupName string) string {
	return fmt.Sprintf(`

	resource "stackit_object_storage_project" "example" {
		project_id         = "%s"
	}

	resource "stackit_object_storage_credentials_group" "example" {
		object_storage_project_id = stackit_object_storage_project.example.id
		name	   				  = "%s"
	}

	resource "stackit_object_storage_credential" "example" {
		object_storage_project_id 	= stackit_object_storage_project.example.id
		credentials_group_id 		= stackit_object_storage_credentials_group.example.id
	}
	  `,
		common.GetAcceptanceTestsProjectID(),
		groupName,
	)
}
