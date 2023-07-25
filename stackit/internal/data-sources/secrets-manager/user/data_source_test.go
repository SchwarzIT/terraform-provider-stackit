package user_test

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
				Config: config(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttrPair("stackit_secrets_manager_user.example", "username", "data.stackit_secrets_manager_user.example", "username"),
					resource.TestCheckTypeSetElemAttrPair("stackit_secrets_manager_user.example", "id", "data.stackit_secrets_manager_user.example", "id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_secrets_manager_user.example", "instance_id", "data.stackit_secrets_manager_user.example", "instance_id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_secrets_manager_user.example", "writable", "data.stackit_secrets_manager_user.example", "writable"),
					resource.TestCheckTypeSetElemAttrPair("stackit_secrets_manager_user.example", "description", "data.stackit_secrets_manager_user.example", "description"),
				),
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

resource "stackit_secrets_manager_user" "example" {
	project_id         = stackit_secrets_manager_instance.example.project_id
	instance_id        = stackit_secrets_manager_instance.example.id
	description        = "test"
	writable           = true
}

data "stackit_secrets_manager_user" "example" {
	project_id         = stackit_secrets_manager_instance.example.project_id
	instance_id        = stackit_secrets_manager_instance.example.id
	username           = stackit_secrets_manager_user.example.username
}
	  `,
		common.GetAcceptanceTestsProjectID(),
		name,
	)
}
