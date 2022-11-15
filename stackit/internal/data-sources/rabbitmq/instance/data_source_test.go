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

const run_this_test = true

func TestAcc_RabbitMQInstance(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}

	name := "odjtest-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)
	plan := "stackit-rabbitmq-single-small"
	planID := "4bc417ff-cb98-4064-bb56-8a2654120768"
	version := "3.7"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(name, version, plan),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.stackit_rabbitmq_instance.example", "name", name),
					resource.TestCheckResourceAttr("data.stackit_rabbitmq_instance.example", "project_id", common.ACC_TEST_PROJECT_ID),
					resource.TestCheckResourceAttr("data.stackit_rabbitmq_instance.example", "version", version),
					resource.TestCheckResourceAttr("data.stackit_rabbitmq_instance.example", "plan", plan),
					resource.TestCheckResourceAttr("data.stackit_rabbitmq_instance.example", "plan_id", planID),
					resource.TestCheckResourceAttrSet("data.stackit_rabbitmq_instance.example", "id"),
					resource.TestCheckResourceAttrSet("data.stackit_rabbitmq_instance.example", "dashboard_url"),
					resource.TestCheckResourceAttrSet("data.stackit_rabbitmq_instance.example", "cf_guid"),
					resource.TestCheckResourceAttrSet("data.stackit_rabbitmq_instance.example", "cf_space_guid"),
					resource.TestCheckResourceAttrSet("data.stackit_rabbitmq_instance.example", "cf_organization_guid"),
					resource.TestCheckTypeSetElemAttrPair("stackit_rabbitmq_instance.example", "id", "data.stackit_rabbitmq_instance.example", "id"),
				),
			},
		},
	})
}

func config(name, version, plan string) string {
	return fmt.Sprintf(`
	resource "stackit_rabbitmq_instance" "example" {
		name       = "%s"
		project_id = "%s"
		version    = "%s"
		plan       = "%s"
	  }


	
	  data "stackit_rabbitmq_instance" "example" {
		depends_on = [stackit_rabbitmq_instance.example]
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
