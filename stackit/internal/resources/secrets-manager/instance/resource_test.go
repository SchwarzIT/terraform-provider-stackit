package instance_test

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

const run_this_test = true

func TestAcc_SecretsManagerInstance(t *testing.T) {
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
					resource.TestCheckResourceAttr("stackit_secrets_manager_instance.example", "name", name),
					resource.TestCheckResourceAttrSet("stackit_secrets_manager_instance.example", "id"),
					resource.TestCheckResourceAttrSet("stackit_secrets_manager_instance.example", "frontend_url"),
					resource.TestCheckResourceAttrSet("stackit_secrets_manager_instance.example", "api_url"),
				),
			},
			{
				Config: config2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_secrets_manager_instance.example", "acl.#", "2"),
				),
			},
			{
				Config: config3(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_secrets_manager_instance.example", "acl.#", "1"),
				),
			},
			// test import
			{
				ResourceName: "stackit_secrets_manager_instance.example",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r, ok := s.RootModule().Resources["stackit_secrets_manager_instance.example"]
					if !ok {
						return "", errors.New("couldn't find resource stackit_secrets_manager_instance.example")
					}
					id, ok := r.Primary.Attributes["id"]
					if !ok {
						return "", errors.New("couldn't find attribute id")
					}

					return fmt.Sprintf("%s,%s", common.GetAcceptanceTestsProjectID(), id), nil
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func config(name string) string {
	return fmt.Sprintf(`
resource "stackit_secrets_manager_instance" "example" {
	project_id         = "%s"
	name               = "%s"
}
	  `,
		common.GetAcceptanceTestsProjectID(),
		name,
	)
}

func config2(name string) string {
	return fmt.Sprintf(`
resource "stackit_secrets_manager_instance" "example" {
	project_id         = "%s"
	name               = "%s"
	acl                = ["193.148.160.0/19","45.129.40.1/21"]
}
	  `,
		common.GetAcceptanceTestsProjectID(),
		name,
	)
}

func config3(name string) string {
	return fmt.Sprintf(`
resource "stackit_secrets_manager_instance" "example" {
	project_id         = "%s"
	name               = "%s"
	acl                = ["193.148.160.0/19"]
}
	  `,
		common.GetAcceptanceTestsProjectID(),
		name,
	)
}
