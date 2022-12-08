package instance_test

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

func TestAcc_PostgresDBInstance(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}

	name := "odjtest-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)
	plan := "stackit-postgresql-single-small"
	planID := "d5752507-13d1-48ee-8ef1-cd6537bd00a4"
	version := "11"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(name, version, plan),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.stackit_postgres_instance.example", "name", name),
					resource.TestCheckResourceAttr("data.stackit_postgres_instance.example", "project_id", common.ACC_TEST_PROJECT_ID),
					resource.TestCheckResourceAttr("data.stackit_postgres_instance.example", "version", version),
					resource.TestCheckResourceAttr("data.stackit_postgres_instance.example", "plan", plan),
					resource.TestCheckResourceAttr("data.stackit_postgres_instance.example", "plan_id", planID),
					resource.TestCheckResourceAttrSet("data.stackit_postgres_instance.example", "id"),
					resource.TestCheckResourceAttrSet("data.stackit_postgres_instance.example", "dashboard_url"),
					resource.TestCheckResourceAttrSet("data.stackit_postgres_instance.example", "cf_guid"),
					resource.TestCheckResourceAttrSet("data.stackit_postgres_instance.example", "cf_space_guid"),
					resource.TestCheckResourceAttrSet("data.stackit_postgres_instance.example", "cf_organization_guid"),
					resource.TestCheckTypeSetElemAttrPair("stackit_postgres_instance.example", "id", "data.stackit_postgres_instance.example", "id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_postgres_instance.example", "dashboard_url", "data.stackit_postgres_instance.example", "dashboard_url"),
					resource.TestCheckTypeSetElemAttrPair("stackit_postgres_instance.example", "cf_space_guid", "data.stackit_postgres_instance.example", "cf_space_guid"),
					resource.TestCheckTypeSetElemAttrPair("stackit_postgres_instance.example", "cf_organization_guid", "data.stackit_postgres_instance.example", "cf_organization_guid"),
				),
			},
		},
	})
}

func config(name, version, plan string) string {
	return fmt.Sprintf(`
	resource "stackit_postgres_instance" "example" {
		name       = "%s"
		project_id = "%s"
		version    = "%s"
		plan       = "%s"
	  }


	
	  data "stackit_postgres_instance" "example" {
		depends_on = [stackit_postgres_instance.example]
		name       = "%s"
		project_id = "%s"
	  }

	`,
		name,
		common.ACC_TEST_PROJECT_ID,
		version,
		plan,
		name,
		common.ACC_TEST_PROJECT_ID,
	)
}