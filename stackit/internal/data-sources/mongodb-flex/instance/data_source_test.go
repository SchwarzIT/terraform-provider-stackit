package instance_test

import (
	"fmt"
	"testing"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/mongodb-flex/instance/defaults"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const run_this_test = false
const default_machine_type = "1.1"

func TestAcc_MongoDBFlexInstance(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.stackit_mongodb_flex_instance.example", "name", name),
					resource.TestCheckResourceAttr("data.stackit_mongodb_flex_instance.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttr("data.stackit_mongodb_flex_instance.example", "machine_type", default_machine_type),
					resource.TestCheckResourceAttr("data.stackit_mongodb_flex_instance.example", "replicas", fmt.Sprint(defaults.Replicas)),
					resource.TestCheckResourceAttr("data.stackit_mongodb_flex_instance.example", "storage.class", defaults.Storage_class),
					resource.TestCheckResourceAttr("data.stackit_mongodb_flex_instance.example", "storage.size", fmt.Sprint(defaults.Storage_size)),
					resource.TestCheckResourceAttrSet("data.stackit_mongodb_flex_instance.example", "id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mongodb_flex_instance.example", "id", "data.stackit_mongodb_flex_instance.example", "id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mongodb_flex_instance.example", "machine_type", "data.stackit_mongodb_flex_instance.example", "machine_type"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mongodb_flex_instance.example", "replicas", "data.stackit_mongodb_flex_instance.example", "replicas"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mongodb_flex_instance.example", "storage.class", "data.stackit_mongodb_flex_instance.example", "storage.class"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mongodb_flex_instance.example", "storage.size", "data.stackit_mongodb_flex_instance.example", "storage.size"),
				),
			},
		},
	})
}

func config(name string) string {
	return fmt.Sprintf(`
	resource "stackit_mongodb_flex_instance" "example" {
		name         = "%s"
		project_id   = "%s"
		machine_type = "%s"
	} 

	data "stackit_mongodb_flex_instance" "example" {
		depends_on = [stackit_mongodb_flex_instance.example]
		name       = "%s"
		project_id = "%s"
	  }

	`,
		name,
		common.GetAcceptanceTestsProjectID(),
		default_machine_type,
		name,
		common.GetAcceptanceTestsProjectID(),
	)
}
