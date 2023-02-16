package user_test

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

func TestAcc_PostgresFlexUser(t *testing.T) {
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
				Config: config(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.stackit_postgres_flex_user.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckTypeSetElemAttrPair("stackit_postgres_flex_user.example", "id", "data.stackit_postgres_flex_user.example", "id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_postgres_flex_user.example", "instance_id", "data.stackit_postgres_flex_user.example", "instance_id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_postgres_flex_user.example", "username", "data.stackit_postgres_flex_user.example", "username"),
					resource.TestCheckTypeSetElemAttrPair("stackit_postgres_flex_user.example", "host", "data.stackit_postgres_flex_user.example", "host"),
					resource.TestCheckTypeSetElemAttrPair("stackit_postgres_flex_user.example", "port", "data.stackit_postgres_flex_user.example", "port"),
				),
			},
		},
	})
}

func config(name string) string {
	return fmt.Sprintf(`
	resource "stackit_postgres_flex_instance" "example" {
		name         = "%s"
		project_id   = "%s"
		machine_type = "c1.2"
	} 
	resource "stackit_postgres_flex_user" "example" {
		project_id   = "%s"
		instance_id  = stackit_postgres_flex_instance.example.id
	}  

	data "stackit_postgres_flex_user" "example" {
		id         = stackit_postgres_flex_user.example.id
		project_id = "%s"
		instance_id = stackit_postgres_flex_instance.example.id
	}

	`,
		name,
		common.GetAcceptanceTestsProjectID(),
		common.GetAcceptanceTestsProjectID(),
		common.GetAcceptanceTestsProjectID(),
	)
}
