package instance_test

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

func TestAcc_RedisJob(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}

	projectId := common.ACC_TEST_PROJECT_ID
	if val, exists := os.LookupEnv("STACKIT_TEST_PROJECT_ID"); exists {
		projectId = val
	}
	name := "odjtest-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)
	plan1 := "stackit-redis-single-small"
	plan2 := "stackit-redis-single-medium"
	planId := "09876364-e1ba-49ec-845c-e8ac45f84921"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(projectId, name, plan1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_redis_instance.example", "name", name),
					resource.TestCheckResourceAttr("stackit_redis_instance.example", "project_id", projectId),
					resource.TestCheckResourceAttr("stackit_redis_instance.example", "version", "6.0"),
					resource.TestCheckResourceAttr("stackit_redis_instance.example", "plan", plan1),
					resource.TestCheckResourceAttr("stackit_redis_instance.example", "plan_id", planId),
					resource.TestCheckResourceAttrSet("stackit_redis_instance.example", "id"),
					resource.TestCheckResourceAttrSet("stackit_redis_instance.example", "plan_id"),
					resource.TestCheckResourceAttrSet("stackit_redis_instance.example", "dashboard_url"),
					resource.TestCheckResourceAttrSet("stackit_redis_instance.example", "cf_guid"),
					resource.TestCheckResourceAttrSet("stackit_redis_instance.example", "cf_space_guid"),
					resource.TestCheckResourceAttrSet("stackit_redis_instance.example", "cf_organization_guid"),
				),
			},
			// check update plan
			{
				Config: config(projectId, name, plan2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_redis_instance.example", "name", name),
					resource.TestCheckResourceAttr("stackit_redis_instance.example", "project_id", projectId),
					resource.TestCheckResourceAttr("stackit_redis_instance.example", "version", "6.0"),
					resource.TestCheckResourceAttr("stackit_redis_instance.example", "plan", plan2),
					resource.TestCheckResourceAttrSet("stackit_redis_instance.example", "id"),
					resource.TestCheckResourceAttrSet("stackit_redis_instance.example", "plan_id"),
					resource.TestCheckResourceAttrSet("stackit_redis_instance.example", "dashboard_url"),
					resource.TestCheckResourceAttrSet("stackit_redis_instance.example", "cf_guid"),
					resource.TestCheckResourceAttrSet("stackit_redis_instance.example", "cf_space_guid"),
					resource.TestCheckResourceAttrSet("stackit_redis_instance.example", "cf_organization_guid"),
				),
			},
			// test import
			{
				ResourceName: "stackit_redis_instance.example",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r, ok := s.RootModule().Resources["stackit_redis_instance.example"]
					if !ok {
						return "", errors.New("couldn't find resource stackit_redis_instance.example")
					}
					id, ok := r.Primary.Attributes["id"]
					if !ok {
						return "", errors.New("couldn't find attribute id")
					}

					return fmt.Sprintf("%s,%s", projectId, id), nil
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func config(projectId, name, plan string) string {
	return fmt.Sprintf(`
	resource "stackit_redis_instance" "example" {
		name       = "%s"
		project_id = "%s"
		version    = "6.0"
		plan       = "%s"
	  }
	  
	  `,
		name,
		projectId,
		plan,
	)
}
