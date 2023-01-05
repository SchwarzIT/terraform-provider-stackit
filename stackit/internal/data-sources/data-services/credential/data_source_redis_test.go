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

const redis_cred_run_this_test = false

func TestAcc_DataSourceRedisCredentialJob(t *testing.T) {
	if !common.ShouldAccTestRun(redis_cred_run_this_test) {
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
				Config: configCredRedis(projectID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.stackit_redis_credential.example", "project_id", projectID),
					resource.TestCheckTypeSetElemAttrPair("stackit_redis_credential.example", "project_id", "data.stackit_redis_credential.example", "project_id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_redis_credential.example", "instance_id", "data.stackit_redis_credential.example", "instance_id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_redis_credential.example", "id", "data.stackit_redis_credential.example", "id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_redis_credential.example", "host", "data.stackit_redis_credential.example", "host"),
					resource.TestCheckTypeSetElemAttrPair("stackit_redis_credential.example", "username", "data.stackit_redis_credential.example", "username"),
					resource.TestCheckTypeSetElemAttrPair("stackit_redis_credential.example", "password", "data.stackit_redis_credential.example", "password"),
					resource.TestCheckTypeSetElemAttrPair("stackit_redis_credential.example", "port", "data.stackit_redis_credential.example", "port"),
					resource.TestCheckTypeSetElemAttrPair("stackit_redis_credential.example", "syslog_drain_url", "data.stackit_redis_credential.example", "syslog_drain_url"),
					resource.TestCheckTypeSetElemAttrPair("stackit_redis_credential.example", "route_service_url", "data.stackit_redis_credential.example", "route_service_url"),
					resource.TestCheckTypeSetElemAttrPair("stackit_redis_credential.example", "uri", "data.stackit_redis_credential.example", "uri"),
				),
			},
		},
	})
}

func configCredRedis(projectID, name string) string {
	return fmt.Sprintf(`
	resource "stackit_redis_instance" "example" {
		name       = "%s"
		project_id = "%s"
	  }
	  
	resource "stackit_redis_credential" "example" {
		project_id = "%s"
		instance_id = stackit_redis_instance.example.id
	}

	data "stackit_redis_credential" "example" {
		depends_on = [stackit_redis_credential.example]
		project_id = "%s"
		instance_id = stackit_redis_instance.example.id
		id = stackit_redis_credential.example.id
	}
	`,
		name,
		projectID,
		projectID,
		projectID,
	)
}
