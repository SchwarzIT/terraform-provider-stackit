package credential_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const postgres_cred_run_this_test = false

func TestAcc_DataSourcePostgresCredentialJob(t *testing.T) {
	if !common.ShouldAccTestRun(postgres_cred_run_this_test) {
		t.Skip()
		return
	}

	name := "odjtest-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)
	projectID := common.GetAcceptanceTestsProjectID()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: configCredPostgres(projectID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.stackit_postgres_credential.example", "project_id", projectID),
					resource.TestCheckTypeSetElemAttrPair("stackit_postgres_credential.example", "project_id", "data.stackit_postgres_credential.example", "project_id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_postgres_credential.example", "instance_id", "data.stackit_postgres_credential.example", "instance_id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_postgres_credential.example", "id", "data.stackit_postgres_credential.example", "id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_postgres_credential.example", "host", "data.stackit_postgres_credential.example", "host"),
					resource.TestCheckTypeSetElemAttrPair("stackit_postgres_credential.example", "username", "data.stackit_postgres_credential.example", "username"),
					resource.TestCheckTypeSetElemAttrPair("stackit_postgres_credential.example", "password", "data.stackit_postgres_credential.example", "password"),
					resource.TestCheckTypeSetElemAttrPair("stackit_postgres_credential.example", "port", "data.stackit_postgres_credential.example", "port"),
					resource.TestCheckTypeSetElemAttrPair("stackit_postgres_credential.example", "syslog_drain_url", "data.stackit_postgres_credential.example", "syslog_drain_url"),
					resource.TestCheckTypeSetElemAttrPair("stackit_postgres_credential.example", "route_service_url", "data.stackit_postgres_credential.example", "route_service_url"),
					resource.TestCheckTypeSetElemAttrPair("stackit_postgres_credential.example", "uri", "data.stackit_postgres_credential.example", "uri"),
				),
			},
		},
	})
}

func configCredPostgres(projectID, name string) string {
	return fmt.Sprintf(`
	resource "stackit_postgres_instance" "example" {
		name       = "%s"
		project_id = "%s"
	  }
	  
	resource "stackit_postgres_credential" "example" {
		project_id = "%s"
		instance_id = stackit_postgres_instance.example.id
	}

	data "stackit_postgres_credential" "example" {
		depends_on = [stackit_postgres_credential.example]
		project_id = "%s"
		instance_id = stackit_postgres_instance.example.id
		id = stackit_postgres_credential.example.id
	}
	`,
		name,
		projectID,
		projectID,
		projectID,
	)
}
