package credential_test

import (
	"errors"
	"fmt"
	"os"
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

func TestAcc_MariaDbCredentialResource(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}

	projectId := common.ACC_TEST_PROJECT_ID
	if val, exists := os.LookupEnv("STACKIT_TEST_PROJECT_ID"); exists {
		projectId = val
	}
	name := "odjtest-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(projectId, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_mariadb_credential.example", "project_id", projectId),
					resource.TestCheckResourceAttrSet("stackit_mariadb_credential.example", "instance_id"),
					resource.TestCheckResourceAttrSet("stackit_mariadb_credential.example", "id"),
					resource.TestCheckResourceAttrSet("stackit_mariadb_credential.example", "host"),
					resource.TestCheckResourceAttrSet("stackit_mariadb_credential.example", "username"),
					resource.TestCheckResourceAttrSet("stackit_mariadb_credential.example", "password"),
					resource.TestCheckResourceAttrSet("stackit_mariadb_credential.example", "port"),
					resource.TestCheckResourceAttrSet("stackit_mariadb_credential.example", "uri"),
				),
			},
			// test import
			{
				ResourceName: "stackit_mariadb_credential.example",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r, ok := s.RootModule().Resources["stackit_mariadb_credential.example"]
					if !ok {
						return "", errors.New("couldn't find resource stackit_mariadb_instance.example")
					}
					id, ok := r.Primary.Attributes["id"]
					if !ok {
						return "", errors.New("couldn't find attribute id")
					}
					iid, ok := r.Primary.Attributes["instance_id"]
					if !ok {
						return "", errors.New("couldn't find attribute id")
					}

					return fmt.Sprintf("%s,%s,%s", projectId, iid, id), nil
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func config(project_id, name string) string {
	return fmt.Sprintf(`
	resource "stackit_mariadb_instance" "example" {
		name       = "%s"
		project_id = "%s"
	}
	  
	resource "stackit_mariadb_credential" "example" {
		project_id = "%s"
		instance_id = stackit_mariadb_instance.example.id
	}
	`,
		name,
		project_id,
		project_id,
	)
}
