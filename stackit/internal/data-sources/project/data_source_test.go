package project_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const run_this_test = true

const schwarz_container_id = "schwarz-it-kg-WJACUK1"

func TestAcc_ProjectDataSource(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}

	billing, ok := os.LookupEnv("ACC_TEST_BILLING_REF")
	if !ok {
		t.Skip("Skipping TestAcc_Project: ACC_TEST_BILLING_REF not specified")
	}

	user, ok := os.LookupEnv("ACC_TEST_USER_EMAIL")
	if !ok {
		t.Skip("Skipping TestAcc_Project: ACC_TEST_USER_EMAIL not specified")
	}

	name := "ODJ AccTest " + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(name, user, billing),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttrPair("stackit_project.example", "name", "data.stackit_project.ex1", "name"),
					resource.TestCheckTypeSetElemAttrPair("stackit_project.example", "billing_ref", "data.stackit_project.ex1", "billing_ref"),
					resource.TestCheckTypeSetElemAttrPair("stackit_project.example", "container_id", "data.stackit_project.ex1", "container_id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_project.example", "parent_container_id", "data.stackit_project.ex1", "parent_container_id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_project.example", "labels", "data.stackit_project.ex1", "labels"),
				),
			},
		},
	})
}

func config(name, owner, billing string) string {
	return fmt.Sprintf(`
	resource "stackit_project" "example" {
		name        = "%s"
		billing_ref = "%s"
		owner_email = "%s"
		parent_container_id = "%s"
		labels = {}
	}

	data "stackit_project" "ex1" {
		depends_on   = [stackit_project.example]
		container_id = stackit_project.example.container_id
	}

	`,
		name, billing, owner, schwarz_container_id,
	)
}
