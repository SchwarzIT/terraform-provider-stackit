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

const run_this_test = false

func TestAcc_PostgresDBJob(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}

	name := "odjtest-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)
	plan1 := "stackit-postgresql-single-small"
	plan2 := "stackit-postgresql-single-medium"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(name, plan1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_postgres_instance.example", "name", name),
					resource.TestCheckResourceAttr("stackit_postgres_instance.example", "project_id", common.ACC_TEST_PROJECT_ID),
					resource.TestCheckResourceAttr("stackit_postgres_instance.example", "version", "13"),
					resource.TestCheckResourceAttr("stackit_postgres_instance.example", "plan", plan1),
					resource.TestCheckResourceAttr("stackit_postgres_instance.example", "plan_id", "215a5d20-5913-4e0f-8e1e-316e75fa09c5"),
					resource.TestCheckResourceAttrSet("stackit_postgres_instance.example", "id"),
					resource.TestCheckResourceAttrSet("stackit_postgres_instance.example", "plan_id"),
					resource.TestCheckResourceAttrSet("stackit_postgres_instance.example", "dashboard_url"),
					resource.TestCheckResourceAttrSet("stackit_postgres_instance.example", "cf_guid"),
					resource.TestCheckResourceAttrSet("stackit_postgres_instance.example", "cf_space_guid"),
					resource.TestCheckResourceAttrSet("stackit_postgres_instance.example", "cf_organization_guid"),
				),
			},
			// check update plan
			{
				Config: config(name, plan2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_postgres_instance.example", "name", name),
					resource.TestCheckResourceAttr("stackit_postgres_instance.example", "project_id", common.ACC_TEST_PROJECT_ID),
					resource.TestCheckResourceAttr("stackit_postgres_instance.example", "version", "6.0"),
					resource.TestCheckResourceAttr("stackit_postgres_instance.example", "plan", plan2),
					resource.TestCheckResourceAttrSet("stackit_postgres_instance.example", "id"),
					resource.TestCheckResourceAttrSet("stackit_postgres_instance.example", "plan_id"),
					resource.TestCheckResourceAttrSet("stackit_postgres_instance.example", "dashboard_url"),
					resource.TestCheckResourceAttrSet("stackit_postgres_instance.example", "cf_guid"),
					resource.TestCheckResourceAttrSet("stackit_postgres_instance.example", "cf_space_guid"),
					resource.TestCheckResourceAttrSet("stackit_postgres_instance.example", "cf_organization_guid"),
				),
			},
			// test import
			{
				ResourceName: "stackit_postgres_instance.example",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r, ok := s.RootModule().Resources["stackit_postgres_instance.example"]
					if !ok {
						return "", errors.New("couldn't find resource stackit_postgres_instance.example")
					}
					id, ok := r.Primary.Attributes["id"]
					if !ok {
						return "", errors.New("couldn't find attribute id")
					}

					return fmt.Sprintf("%s,%s", common.ACC_TEST_PROJECT_ID, id), nil
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func config(name, plan string) string {
	return fmt.Sprintf(`
	resource "stackit_postgres_instance" "example" {
		name       = "%s"
		project_id = "%s"
		version    = "13"
		plan       = "%s"
	  }
	  
	  `,
		name,
		common.ACC_TEST_PROJECT_ID,
		plan,
	)
}
