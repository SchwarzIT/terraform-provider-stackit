package credential_test

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

func TestAcc_ResourceArgusCredentialJob(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}

	name := "argus" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)
	projectID := common.GetAcceptanceTestsProjectID()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(projectID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_argus_credential.example", "project_id", projectID),
					resource.TestCheckResourceAttrSet("stackit_argus_credential.example", "instance_id"),
					resource.TestCheckResourceAttrSet("stackit_argus_credential.example", "id"),
					resource.TestCheckResourceAttrSet("stackit_argus_credential.example", "username"),
					resource.TestCheckResourceAttrSet("stackit_argus_credential.example", "password"),
				),
			},
			// test import
			{
				ResourceName: "stackit_argus_credential.example",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r, ok := s.RootModule().Resources["stackit_argus_credential.example"]
					if !ok {
						return "", errors.New("couldn't find resource stackit_argus_credential.example")
					}
					username, ok := r.Primary.Attributes["username"]
					if !ok {
						return "", errors.New("couldn't find attribute id")
					}
					iid, ok := r.Primary.Attributes["instance_id"]
					if !ok {
						return "", errors.New("couldn't find attribute id")
					}

					return fmt.Sprintf("%s,%s,%s", common.GetAcceptanceTestsProjectID(), iid, username), nil
				},
				ImportStateVerifyIgnore: []string{"password"},
				ImportState:             true,
				ImportStateVerify:       true,
			},
		},
	})
}

func config(project_id, name string) string {
	return fmt.Sprintf(`
	resource "stackit_argus_instance" "example" {
		project_id = "%s"
		name       = "%s"
		plan       = "Monitoring-Medium-EU01"
	}
	  
	resource "stackit_argus_credential" "example" {
		project_id = "%s"
		instance_id = stackit_argus_instance.example.id
	}
	`,
		project_id,
		name,
		project_id,
	)
}
