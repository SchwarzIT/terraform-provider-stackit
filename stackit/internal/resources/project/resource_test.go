package project_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/consts"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const run_this_test = false

func TestAcc_Project(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
	}

	billingRef, ok := os.LookupEnv("ACC_TEST_BILLING_REF")
	if !ok {
		t.Skip("Skipping TestAcc_Project: ACC_TEST_BILLING_REF not specified")
	}

	user, ok := os.LookupEnv("ACC_TEST_USER_EMAIL")
	if !ok {
		t.Skip("Skipping TestAcc_Project: ACC_TEST_USER_EMAIL not specified")
	}

	name := "ODJ AccTest " + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	newName := "ODJ AccTest " + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: config(name, billingRef, user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("stackit_project.example", "id"),
					resource.TestCheckResourceAttr("stackit_project.example", "name", name),
					resource.TestCheckResourceAttr("stackit_project.example", "billing_ref", billingRef),
				),
			},
			// rename
			{
				Config: config(newName, billingRef, user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("stackit_project.example", "id"),
					resource.TestCheckResourceAttr("stackit_project.example", "name", newName),
					resource.TestCheckResourceAttr("stackit_project.example", "billing_ref", billingRef),
				),
			},
			// enabled services
			{
				Config: config2(newName, billingRef, user, true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("stackit_project.example", "id"),
					resource.TestCheckResourceAttr("stackit_project.example", "name", newName),
					resource.TestCheckResourceAttr("stackit_project.example", "billing_ref", billingRef),
					resource.TestCheckResourceAttr("stackit_project.example", "enable_kubernetes", "true"),
					resource.TestCheckResourceAttr("stackit_project.example", "enable_object_storage", "false"),
				),
			},
			{
				Config: config2(newName, billingRef, user, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("stackit_project.example", "id"),
					resource.TestCheckResourceAttr("stackit_project.example", "name", newName),
					resource.TestCheckResourceAttr("stackit_project.example", "billing_ref", billingRef),
					resource.TestCheckResourceAttr("stackit_project.example", "enable_kubernetes", "false"),
					resource.TestCheckResourceAttr("stackit_project.example", "enable_object_storage", "true"),
				),
			},
			// back to default
			{
				Config: config(newName, billingRef, user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("stackit_project.example", "id"),
					resource.TestCheckResourceAttr("stackit_project.example", "name", newName),
					resource.TestCheckResourceAttr("stackit_project.example", "billing_ref", billingRef),
				),
			},
			// test import
			{
				ResourceName:            "stackit_project.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"owner_email"},
			},
		},
	})
}

func config(name, billingRef, user string) string {
	return fmt.Sprintf(`
	resource "stackit_project" "example" {
		name        = "%s"
		billing_ref = "%s"
		owner_email = "%s"
		parent_container_id = "%s"
	}
	`,
		name,
		billingRef,
		user,
		consts.SCHWARZ_CONTAINER_ID,
	)
}

func config2(name, billingRef, user string, enableKubernetes, enableObjectStorage bool) string {
	return fmt.Sprintf(`
	resource "stackit_project" "example" {
		name        = "%s"
		billing_ref = "%s"
		owner_email = "%s"
		enable_kubernetes = %v
		enable_object_storage = %v
		parent_container_id = "%s"
	}
	`,
		name,
		billingRef,
		user,
		enableKubernetes,
		enableObjectStorage,
		consts.SCHWARZ_CONTAINER_ID,
	)
}
