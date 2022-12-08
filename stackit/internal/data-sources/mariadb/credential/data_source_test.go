package credential_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const run_this_test = true

func TestAcc_MariaDbCredentialDataSource(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.stackit_mariadb_credential.example", "project_id", projectId),
					resource.TestCheckTypeSetElemAttrPair("stackit_mariadb_credential.example", "project_id", "data.stackit_mariadb_credential.example", "project_id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mariadb_credential.example", "instance_id", "data.stackit_mariadb_credential.example", "instance_id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mariadb_credential.example", "id", "data.stackit_mariadb_credential.example", "id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mariadb_credential.example", "ca_cert", "data.stackit_mariadb_credential.example", "ca_cert"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mariadb_credential.example", "host", "data.stackit_mariadb_credential.example", "host"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mariadb_credential.example", "username", "data.stackit_mariadb_credential.example", "username"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mariadb_credential.example", "password", "data.stackit_mariadb_credential.example", "password"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mariadb_credential.example", "port", "data.stackit_mariadb_credential.example", "port"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mariadb_credential.example", "syslog_drain_url", "data.stackit_mariadb_credential.example", "syslog_drain_url"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mariadb_credential.example", "route_service_url", "data.stackit_mariadb_credential.example", "route_service_url"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mariadb_credential.example", "schema", "data.stackit_mariadb_credential.example", "schema"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mariadb_credential.example", "uri", "data.stackit_mariadb_credential.example", "uri"),
				),
			},
		},
	})
}

func config(projectId, name string) string {
	return fmt.Sprintf(`
	resource "stackit_mariadb_instance" "example" {
		name       = "%s"
		project_id = "%s"
	  }
	  
	resource "stackit_mariadb_credential" "example" {
		project_id = "%s"
		instance_id = stackit_mariadb_instance.example.id
	}

	data "stackit_mariadb_credential" "example" {
		project_id = "%s"
		instance_id = stackit_mariadb_instance.example.id
		id = stackit_mariadb_credential.example.id
	}
	`,
		name,
		projectId,
		projectId,
		projectId,
	)
}