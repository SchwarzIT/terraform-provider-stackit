package credential_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"testing"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const run_this_test = false

func TestAcc_ElasticSearchInstance(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}

	iName := "odjtest-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)
	iVersion := "7"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(common.ACC_TEST_PROJECT_ID, iName, iVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_elasticsearch_instance.example", "project_id", common.ACC_TEST_PROJECT_ID),
					resource.TestCheckResourceAttrSet("stackit_elasticsearch_credential.example", "id"),
				),
			},
		},
	})
}

func config(project_id, instance_name, instance_version string) string {
	return fmt.Sprintf(`
	resource "stackit_elasticsearch_instance" "example" {
		project_id = "%s"
		name       = "%s"
		version    = "%s"
	}

	resource "stackit_elasticsearch_credential" "example" {
		project_id  = "%s"
		instance_id = [stackit_elasticsearch_instance.example.ID]
	}

	data "stackit_elasticsearch_credential" "example" {
		depends_on  = [stackit_elasticsearch_credential.example]
		project_id  = "%s"
		instance_id = stackit_elasticsearch_instance.example.ID
	}

	`,
		project_id,
		instance_name,
		instance_version,
		project_id,
		project_id,
	)
}
