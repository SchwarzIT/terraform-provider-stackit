package postgresinstance_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	postgresinstance "github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/postgres-flex/instance"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const run_this_test = false

func TestAcc_PostgresFlexInstance(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}

	name1 := "odjtest-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)
	name2 := "odjtest-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(name1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_postgres_flex_instance.example", "name", name1),
					resource.TestCheckResourceAttr("stackit_postgres_flex_instance.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttr("stackit_postgres_flex_instance.example", "version", "14"),
					resource.TestCheckResourceAttr("stackit_postgres_flex_instance.example", "machine_type", postgresinstance.DefaultMachineType),
					resource.TestCheckResourceAttr("stackit_postgres_flex_instance.example", "replicas", "1"),
					resource.TestCheckResourceAttr("stackit_postgres_flex_instance.example", "storage.class", "premium-perf6-stackit"),
					resource.TestCheckResourceAttr("stackit_postgres_flex_instance.example", "storage.size", fmt.Sprintf("%d", postgresinstance.DefaultStorageSize)),
					resource.TestCheckResourceAttrSet("stackit_postgres_flex_instance.example", "id"),
				),
			},
			// test import
			{
				Config: config(name2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_postgres_flex_instance.example", "name", name2),
					resource.TestCheckResourceAttr("stackit_postgres_flex_instance.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttr("stackit_postgres_flex_instance.example", "version", "14"),
					resource.TestCheckResourceAttr("stackit_postgres_flex_instance.example", "machine_type", postgresinstance.DefaultMachineType),
					resource.TestCheckResourceAttr("stackit_postgres_flex_instance.example", "replicas", "1"),
					resource.TestCheckResourceAttr("stackit_postgres_flex_instance.example", "storage.class", "premium-perf6-stackit"),
					resource.TestCheckResourceAttr("stackit_postgres_flex_instance.example", "storage.size", fmt.Sprintf("%d", postgresinstance.DefaultStorageSize)),
					resource.TestCheckResourceAttrSet("stackit_postgres_flex_instance.example", "id"),
				),
			},
			{
				ResourceName: "stackit_postgres_flex_instance.example",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r, ok := s.RootModule().Resources["stackit_postgres_flex_instance.example"]
					if !ok {
						return "", errors.New("couldn't find resource stackit_postgres_flex_instance.example")
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
	resource "stackit_postgres_flex_instance" "example" {
		name         = "%s"
		project_id   = "%s"
		machine_type = "%s"
		version      = "14"
		replicas     = 1
		storage 	 = {
			class = "premium-perf6-stackit"
			size  = %d
		}
		acl = ["193.148.160.0/19","45.129.40.1/21"]
	}  
	  `,
		name,
		common.GetAcceptanceTestsProjectID(),
		postgresinstance.DefaultMachineType,
		postgresinstance.DefaultStorageSize,
	)
}
