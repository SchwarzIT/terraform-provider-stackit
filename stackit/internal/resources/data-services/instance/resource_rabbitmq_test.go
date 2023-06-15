package instance_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const rabbitmq_inst_run_this_test = false

func TestAcc_ResourceRabbitMQInstanceJob(t *testing.T) {
	if !common.ShouldAccTestRun(rabbitmq_inst_run_this_test) {
		t.Skip()
		return
	}

	name := "odjtest-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)
	plan1 := "stackit-rabbitmq-2.4.10-single"
	planID1 := "7e1f8394-5dd5-40b1-8608-16b4344eb51b"
	plan2 := "stackit-rabbitmq-4.8.50-single"
	planID2 := "c5e2758b-4611-49a7-99a1-2df20aca5616"
	version := "3.10"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: configInstRabbitMQ(name, plan1, version),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_rabbitmq_instance.example", "name", name),
					resource.TestCheckResourceAttr("stackit_rabbitmq_instance.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttr("stackit_rabbitmq_instance.example", "version", version),
					resource.TestCheckResourceAttr("stackit_rabbitmq_instance.example", "plan", plan1),
					resource.TestCheckResourceAttr("stackit_rabbitmq_instance.example", "plan_id", planID1),
					resource.TestCheckResourceAttrSet("stackit_rabbitmq_instance.example", "id"),
					resource.TestCheckResourceAttrSet("stackit_rabbitmq_instance.example", "dashboard_url"),
					resource.TestCheckResourceAttrSet("stackit_rabbitmq_instance.example", "cf_guid"),
					resource.TestCheckResourceAttrSet("stackit_rabbitmq_instance.example", "cf_space_guid"),
				),
			},
			// check update plan
			{
				Config: configInstRabbitMQ(name, plan2, version),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_rabbitmq_instance.example", "name", name),
					resource.TestCheckResourceAttr("stackit_rabbitmq_instance.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttr("stackit_rabbitmq_instance.example", "version", version),
					resource.TestCheckResourceAttr("stackit_rabbitmq_instance.example", "plan", plan2),
					resource.TestCheckResourceAttr("stackit_rabbitmq_instance.example", "plan_id", planID2),
					resource.TestCheckResourceAttrSet("stackit_rabbitmq_instance.example", "id"),
					resource.TestCheckResourceAttrSet("stackit_rabbitmq_instance.example", "dashboard_url"),
					resource.TestCheckResourceAttrSet("stackit_rabbitmq_instance.example", "cf_guid"),
					resource.TestCheckResourceAttrSet("stackit_rabbitmq_instance.example", "cf_space_guid"),
				),
			},
			// test import
			{
				ResourceName: "stackit_rabbitmq_instance.example",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r, ok := s.RootModule().Resources["stackit_rabbitmq_instance.example"]
					if !ok {
						return "", errors.New("couldn't find resource stackit_rabbitmq_instance.example")
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

func configInstRabbitMQ(name, plan, version string) string {
	return fmt.Sprintf(`
	resource "stackit_rabbitmq_instance" "example" {
		name       = "%s"
		project_id = "%s"
		version    = "%s"
		plan       = "%s"
	  }
	  
	  `,
		name,
		common.GetAcceptanceTestsProjectID(),
		version,
		plan,
	)
}
