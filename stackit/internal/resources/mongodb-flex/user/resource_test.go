package user_test

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

const run_this_test = false

func TestAcc_MongoDBFlexUser(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}

	name1 := "odjtest-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(name1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("stackit_mongodb_flex_user.example", "id"),
					resource.TestCheckResourceAttr("stackit_mongodb_flex_user.example", "username", "stackit"),
					resource.TestCheckResourceAttr("stackit_mongodb_flex_user.example", "database", "stackit"),
					resource.TestCheckResourceAttrSet("stackit_mongodb_flex_user.example", "password"),
					resource.TestCheckResourceAttrSet("stackit_mongodb_flex_user.example", "host"),
					resource.TestCheckResourceAttrSet("stackit_mongodb_flex_user.example", "port"),
					resource.TestCheckResourceAttrSet("stackit_mongodb_flex_user.example", "uri"),
				),
			},
			// test import
			{
				ResourceName: "stackit_mongodb_flex_user.example",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r, ok := s.RootModule().Resources["stackit_mongodb_flex_instance.example"]
					if !ok {
						return "", errors.New("couldn't find resource stackit_mongodb_flex_instance.example")
					}
					id, ok := r.Primary.Attributes["id"]
					if !ok {
						return "", errors.New("couldn't find attribute id")
					}
					r2, ok := s.RootModule().Resources["stackit_mongodb_flex_user.example"]
					if !ok {
						return "", errors.New("couldn't find resource stackit_mongodb_flex_user.example")
					}
					id2, ok := r2.Primary.Attributes["id"]
					if !ok {
						return "", errors.New("couldn't find attribute id")
					}

					return fmt.Sprintf("%s,%s,%s", common.GetAcceptanceTestsProjectID(), id, id2), nil
				},
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "uri"},
			},
		},
	})
}

func config(name string) string {
	return fmt.Sprintf(`
	resource "stackit_mongodb_flex_instance" "example" {
	 	name         = "%s"
	 	project_id   = "%s"
	 	machine_type = "1.1"
		acl = ["193.148.160.0/19","45.129.40.1/21"]
	}  
	resource "stackit_mongodb_flex_user" "example" {
		project_id   = "%s"
		instance_id  = stackit_mongodb_flex_instance.example.id
	}  
	  `,
		name,
		common.GetAcceptanceTestsProjectID(),
		common.GetAcceptanceTestsProjectID(),
	)
}
