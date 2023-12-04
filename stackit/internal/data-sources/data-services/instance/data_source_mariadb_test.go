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

const mariadb_inst_run_this_test = false

func TestAcc_DataSourceMariaDBInstanceJob(t *testing.T) {
	if !common.ShouldAccTestRun(mariadb_inst_run_this_test) {
		t.Skip()
		return
	}

	name := "odjtest-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)
	plan := "stackit-mariadb-1.4.10-single"
	planID := "683be856-3587-42de-b1b5-a792ff854f52"
	version := "10.6"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: configMariaDB(name, version, plan),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.stackit_mariadb_instance.example", "name", name),
					resource.TestCheckResourceAttr("data.stackit_mariadb_instance.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttr("data.stackit_mariadb_instance.example", "version", version),
					resource.TestCheckResourceAttr("data.stackit_mariadb_instance.example", "plan", plan),
					resource.TestCheckResourceAttr("data.stackit_mariadb_instance.example", "plan_id", planID),
					resource.TestCheckResourceAttrSet("data.stackit_mariadb_instance.example", "id"),
					resource.TestCheckResourceAttrSet("data.stackit_mariadb_instance.example", "dashboard_url"),
					resource.TestCheckResourceAttrSet("data.stackit_mariadb_instance.example", "cf_guid"),
					resource.TestCheckResourceAttrSet("data.stackit_mariadb_instance.example", "cf_space_guid"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mariadb_instance.example", "id", "data.stackit_mariadb_instance.example", "id"),
				),
			},
		},
	})
}

func configMariaDB(name, version, plan string) string {
	return fmt.Sprintf(`
	resource "stackit_mariadb_instance" "example" {
		name       = "%s"
		project_id = "%s"
		version    = "%s"
		plan       = "%s"
	  }

	
	  data "stackit_mariadb_instance" "example" {
		depends_on = [stackit_mariadb_instance.example]
		name       = "%s"
		project_id = "%s"
	  }

	`,
		name,
		common.GetAcceptanceTestsProjectID(),
		version,
		plan,
		name,
		common.GetAcceptanceTestsProjectID(),
	)
}
