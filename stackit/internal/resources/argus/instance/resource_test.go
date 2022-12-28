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

const run_this_test = false

func TestAcc_ArgusInstances(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}

	name := "e1-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)
	newName := "e2-" + acctest.RandStringFromCharSet(7, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: config(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_argus_instance.example", "name", name),
					resource.TestCheckResourceAttr("stackit_argus_instance.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttrSet("stackit_argus_instance.example", "plan_id"),
					resource.TestCheckResourceAttrSet("stackit_argus_instance.example", "dashboard_url"),
					resource.TestCheckResourceAttrSet("stackit_argus_instance.example", "is_updatable"),
					resource.TestCheckResourceAttrSet("stackit_argus_instance.example", "grafana_url"),
					resource.TestCheckResourceAttrSet("stackit_argus_instance.example", "grafana_initial_admin_password"),
					resource.TestCheckResourceAttrSet("stackit_argus_instance.example", "grafana_initial_admin_user"),
					resource.TestCheckResourceAttrSet("stackit_argus_instance.example", "metrics_url"),
					resource.TestCheckResourceAttrSet("stackit_argus_instance.example", "metrics_push_url"),
					resource.TestCheckResourceAttrSet("stackit_argus_instance.example", "targets_url"),
					resource.TestCheckResourceAttrSet("stackit_argus_instance.example", "alerting_url"),
					resource.TestCheckResourceAttrSet("stackit_argus_instance.example", "logs_url"),
					resource.TestCheckResourceAttrSet("stackit_argus_instance.example", "logs_push_url"),
					resource.TestCheckResourceAttrSet("stackit_argus_instance.example", "jaeger_traces_url"),
					resource.TestCheckResourceAttrSet("stackit_argus_instance.example", "jaeger_ui_url"),
					resource.TestCheckResourceAttrSet("stackit_argus_instance.example", "otlp_traces_url"),
					resource.TestCheckResourceAttrSet("stackit_argus_instance.example", "zipkin_spans_url"),
				),
			},
			// check update
			{
				Config: configExtended(name, "Monitoring-Medium-EU01"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_argus_instance.example", "name", name),
					resource.TestCheckResourceAttr("stackit_argus_instance.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttr("stackit_argus_instance.example", "plan", "Monitoring-Medium-EU01"),
					resource.TestCheckResourceAttr("stackit_argus_instance.example", "grafana.enable_public_access", "true"),
					resource.TestCheckResourceAttr("stackit_argus_instance.example", "metrics.retention_days", "60"),
					resource.TestCheckResourceAttr("stackit_argus_instance.example", "metrics.retention_days_5m_downsampling", "20"),
					resource.TestCheckResourceAttr("stackit_argus_instance.example", "metrics.retention_days_1h_downsampling", "10"),
				),
			},
			// new name
			{
				Config: configExtended(newName, "Monitoring-Medium-EU01"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_argus_instance.example", "name", newName),
					resource.TestCheckResourceAttr("stackit_argus_instance.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttr("stackit_argus_instance.example", "plan", "Monitoring-Medium-EU01"),
				),
			},
			// new plan
			{
				Config: configExtended(newName, "Monitoring-Basic-EU01"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_argus_instance.example", "name", newName),
					resource.TestCheckResourceAttr("stackit_argus_instance.example", "project_id", common.GetAcceptanceTestsProjectID()),
					resource.TestCheckResourceAttr("stackit_argus_instance.example", "plan", "Monitoring-Basic-EU01"),
				),
			},
			// test import
			{
				ResourceName: "stackit_argus_instance.example",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r, ok := s.RootModule().Resources["stackit_argus_instance.example"]
					if !ok {
						return "", errors.New("couldn't find resource stackit_argus_instance.example")
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

func config(name string) string {
	return fmt.Sprintf(`
resource "stackit_argus_instance" "example" {
	project_id = "%s"
	name       = "%s"
	plan       = "Monitoring-Medium-EU01"
}
	  `,
		common.GetAcceptanceTestsProjectID(),
		name,
	)
}

func configExtended(name, plan string) string {
	return fmt.Sprintf(`
resource "stackit_argus_instance" "example" {
	project_id = "%s"
	name       = "%s"
	plan       = "%s"
	grafana	   = {
		enable_public_access = true
	}
	metrics	   = {
		retention_days 				   = 60
		retention_days_5m_downsampling = 20
		retention_days_1h_downsampling = 10
	}
}
	  `,
		common.GetAcceptanceTestsProjectID(),
		name,
		plan,
	)
}
