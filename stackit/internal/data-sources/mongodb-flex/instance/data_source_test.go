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

const run_this_test = false

func TestAcc_MongoDBInstance(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}

	name := "odjtest-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)
	version := "5.0"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(name, version),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.stackit_mongodb_flex_instance.example", "name", name),
					resource.TestCheckResourceAttr("data.stackit_mongodb_flex_instance.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttr("data.stackit_mongodb_flex_instance.example", "version", version),
					resource.TestCheckResourceAttr("data.stackit_mongodb_flex_instance.example", "machine_type", "C1.1"),
					resource.TestCheckResourceAttr("data.stackit_mongodb_flex_instance.example", "replicas", "1"),
					resource.TestCheckResourceAttr("data.stackit_mongodb_flex_instance.example", "storage.class", "premium-perf2-mongodb"),
					resource.TestCheckResourceAttr("data.stackit_mongodb_flex_instance.example", "storage.size", "10"),
					resource.TestCheckResourceAttrSet("data.stackit_mongodb_flex_instance.example", "id"),
					resource.TestCheckResourceAttrSet("data.stackit_mongodb_flex_instance.example", "user.id"),
					resource.TestCheckResourceAttr("data.stackit_mongodb_flex_instance.example", "user.username", "stackit"),
					resource.TestCheckResourceAttr("data.stackit_mongodb_flex_instance.example", "user.database", "stackit"),
					resource.TestCheckResourceAttrSet("data.stackit_mongodb_flex_instance.example", "user.password"),
					resource.TestCheckResourceAttrSet("data.stackit_mongodb_flex_instance.example", "user.hostname"),
					resource.TestCheckResourceAttrSet("data.stackit_mongodb_flex_instance.example", "user.port"),
					resource.TestCheckResourceAttrSet("data.stackit_mongodb_flex_instance.example", "user.uri"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mongodb_flex_instance.example", "id", "data.stackit_mongodb_flex_instance.example", "id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mongodb_flex_instance.example", "machine_type", "data.stackit_mongodb_flex_instance.example", "machine_type"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mongodb_flex_instance.example", "replicas", "data.stackit_mongodb_flex_instance.example", "replicas"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mongodb_flex_instance.example", "storage.class", "data.stackit_mongodb_flex_instance.example", "storage.class"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mongodb_flex_instance.example", "storage.size", "data.stackit_mongodb_flex_instance.example", "storage.size"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mongodb_flex_instance.example", "user.id", "data.stackit_mongodb_flex_instance.example", "user.id"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mongodb_flex_instance.example", "user.username", "data.stackit_mongodb_flex_instance.example", "user.username"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mongodb_flex_instance.example", "user.hostname", "data.stackit_mongodb_flex_instance.example", "user.hostname"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mongodb_flex_instance.example", "user.port", "data.stackit_mongodb_flex_instance.example", "user.port"),
					resource.TestCheckTypeSetElemAttrPair("stackit_mongodb_flex_instance.example", "user.uri", "data.stackit_mongodb_flex_instance.example", "user.uri"),
				),
			},
		},
	})
}

func config(name, version string) string {
	return fmt.Sprintf(`
	resource "stackit_mongodb_flex_instance" "example" {
		name         = "%s"
		project_id   = "%s"
		machine_type = "C1.1"
		version 	 = "%s"
	}


	
	  data "stackit_mongodb_flex_instance" "example" {
		depends_on = [stackit_mongodb_flex_instance.example]
		name       = "%s"
		project_id = "%s"
	  }

	`,
		name,
		common.GetAcceptanceTestsProjectID(),
		version,
		name,
		common.GetAcceptanceTestsProjectID(),
	)
}
