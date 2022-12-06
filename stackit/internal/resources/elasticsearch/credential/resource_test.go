package credential_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const run_this_test = false

func TestAcc_ElasticSearchJob(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}

	const project_id = "some project id"
	const instance_id = "some instance id"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(project_id, instance_id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_elasticsearch_instance.example", "project_id", project_id),
					resource.TestCheckResourceAttr("stackit_elasticsearch_instance.example", "instance_id", instance_id),
					resource.TestCheckResourceAttrSet("stackit_elasticsearch_credential.example", "id"),
				),
			},
			// test import
			{
				ResourceName: "stackit_elasticsearch_instance.example",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r, ok := s.RootModule().Resources["stackit_elasticsearch_instance.example"]
					if !ok {
						return "", errors.New("couldn't find resource stackit_elasticsearch_instance.example")
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

func config(project_id, instance_id string) string {
	return fmt.Sprintf(`
	resource "stackit_elasticsearch_credential" "example" {
		project_id = "%s"
		instance_id = "%s"
	}
	`,
		project_id,
		instance_id,
	)
}
