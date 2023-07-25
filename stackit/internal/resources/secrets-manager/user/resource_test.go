package user_test

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

func TestAcc_SecretsManagerUser(t *testing.T) {
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
				Config: config(name, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_secrets_manager_user.example", "description", "test"),
					resource.TestCheckResourceAttr("stackit_secrets_manager_user.example", "write_enabled", "true"),
					resource.TestCheckResourceAttrSet("stackit_secrets_manager_user.example", "username"),
					resource.TestCheckResourceAttrSet("stackit_secrets_manager_user.example", "password"),
					resource.TestCheckResourceAttrSet("stackit_secrets_manager_user.example", "id"),
				),
			},
			{
				Config: config(name, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_secrets_manager_user.example", "write_enabled", "false"),
				),
			},
			// test import
			{
				ResourceName: "stackit_secrets_manager_user.example",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r, ok := s.RootModule().Resources["stackit_secrets_manager_user.example"]
					if !ok {
						return "", errors.New("couldn't find resource stackit_secrets_manager_user.example")
					}
					id, ok := r.Primary.Attributes["id"]
					if !ok {
						return "", errors.New("couldn't find attribute id")
					}
					iid, ok := r.Primary.Attributes["instance_id"]
					if !ok {
						return "", errors.New("couldn't find attribute instance_id")
					}

					return fmt.Sprintf("%s,%s,%s", common.GetAcceptanceTestsProjectID(), iid, id), nil
				},
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func config(name string, writeable bool) string {
	return fmt.Sprintf(`
	resource "stackit_secrets_manager_instance" "example" {
		project_id         = "%s"
		name               = "%s"
	}

	resource "stackit_secrets_manager_user" "example" {
		project_id         = stackit_secrets_manager_instance.example.project_id
		instance_id        = stackit_secrets_manager_instance.example.id
		description        = "test"
		write_enabled      = %v
	}
	  `,
		common.GetAcceptanceTestsProjectID(),
		name,
		writeable,
	)
}
