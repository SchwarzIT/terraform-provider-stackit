package credential_test

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

const rabbitmq_cred_run_this_test = false

func TestAcc_ResourceRabbitMQCredentialJob(t *testing.T) {
	if !common.ShouldAccTestRun(rabbitmq_cred_run_this_test) {
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
				Config: configCredRabbitMQ(projectID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_rabbitmq_credential.example", "project_id", projectID),
					resource.TestCheckResourceAttrSet("stackit_rabbitmq_credential.example", "instance_id"),
					resource.TestCheckResourceAttrSet("stackit_rabbitmq_credential.example", "id"),
					resource.TestCheckResourceAttrSet("stackit_rabbitmq_credential.example", "host"),
					resource.TestCheckResourceAttrSet("stackit_rabbitmq_credential.example", "password"),
					resource.TestCheckResourceAttrSet("stackit_rabbitmq_credential.example", "port"),
					resource.TestCheckResourceAttrSet("stackit_rabbitmq_credential.example", "uri"),
					resource.TestCheckResourceAttrSet("stackit_rabbitmq_credential.example", "raw_response"),
				),
			},
			// test import
			{
				ResourceName: "stackit_rabbitmq_credential.example",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r, ok := s.RootModule().Resources["stackit_rabbitmq_credential.example"]
					if !ok {
						return "", errors.New("couldn't find resource stackit_rabbitmq_instance.example")
					}
					id, ok := r.Primary.Attributes["id"]
					if !ok {
						return "", errors.New("couldn't find attribute id")
					}
					iid, ok := r.Primary.Attributes["instance_id"]
					if !ok {
						return "", errors.New("couldn't find attribute id")
					}

					return fmt.Sprintf("%s,%s,%s", common.GetAcceptanceTestsProjectID(), iid, id), nil
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func configCredRabbitMQ(project_id, name string) string {
	return fmt.Sprintf(`
	resource "stackit_rabbitmq_instance" "example" {
		name       = "%s"
		project_id = "%s"
	  }
	  
	resource "stackit_rabbitmq_credential" "example" {
		project_id = "%s"
		instance_id = stackit_rabbitmq_instance.example.id
	}
	`,
		name,
		project_id,
		project_id,
	)
}
