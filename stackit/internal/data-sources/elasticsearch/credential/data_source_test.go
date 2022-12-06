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

const run_this_test = true

func TestAcc_ElasticSearchCredential(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}

	name := "odjtest-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(common.ACC_TEST_PROJECT_ID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.stackit_elasticsearch_credential.example", "project_id", common.ACC_TEST_PROJECT_ID),
					resource.TestCheckTypeSetElemAttrPair("stackit_elasticsearch_credential.example", "project_id", "data.stackit_elasticsearch_credential.example", "project_id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_elasticsearch_credential.example", "instance_id", "data.stackit_elasticsearch_credential.example", "instance_id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_elasticsearch_credential.example", "id", "data.stackit_elasticsearch_credential.example", "id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_elasticsearch_credential.example", "ca_cert", "data.stackit_elasticsearch_credential.example", "ca_cert"),
					resource.TestCheckTypeSetElemAttrPair("stackit_elasticsearch_credential.example", "host", "data.stackit_elasticsearch_credential.example", "host"),
					resource.TestCheckTypeSetElemAttrPair("stackit_elasticsearch_credential.example", "username", "data.stackit_elasticsearch_credential.example", "username"),
					resource.TestCheckTypeSetElemAttrPair("stackit_elasticsearch_credential.example", "password", "data.stackit_elasticsearch_credential.example", "password"),
					resource.TestCheckTypeSetElemAttrPair("stackit_elasticsearch_credential.example", "port", "data.stackit_elasticsearch_credential.example", "port"),
					resource.TestCheckTypeSetElemAttrPair("stackit_elasticsearch_credential.example", "syslog_drain_url", "data.stackit_elasticsearch_credential.example", "syslog_drain_url"),
					resource.TestCheckTypeSetElemAttrPair("stackit_elasticsearch_credential.example", "route_service_url", "data.stackit_elasticsearch_credential.example", "route_service_url"),
					resource.TestCheckTypeSetElemAttrPair("stackit_elasticsearch_credential.example", "schema", "data.stackit_elasticsearch_credential.example", "schema"),
					resource.TestCheckTypeSetElemAttrPair("stackit_elasticsearch_credential.example", "uri", "data.stackit_elasticsearch_credential.example", "uri"),
				),
			},
		},
	})
}

func config(project_id, name string) string {
	return fmt.Sprintf(`
	resource "stackit_elasticsearch_instance" "example" {
		name       = "%s"
		project_id = "%s"
		version    = "7"
	  }
	  
	resource "stackit_elasticsearch_credential" "example" {
		project_id = "%s"
		instance_id = stackit_elasticsearch_instance.example.id
	}

	data "stackit_elasticsearch_credential" "example" {
		project_id = "%s"
		instance_id = stackit_elasticsearch_instance.example.id
		id = stackit_elasticsearch_credential.example.id
	}
	`,
		name,
		project_id,
		project_id,
		project_id,
	)
}
