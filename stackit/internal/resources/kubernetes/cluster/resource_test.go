package cluster_test

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

func TestAcc_kubernetes(t *testing.T) {
	if !common.ShouldAccTestRun(run_this_test) {
		t.Skip()
		return
	}
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stackit": providerserver.NewProtocol6WithError(stackit.New("test")()),
		},
		Steps: []resource.TestStep{
			// check minimal configuration
			{
				Config: configMinimal(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "name", name),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "allow_privileged_containers", "true"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.name", "example-np"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.machine_type", "c1.2"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.os_name", "flatcar"),
					resource.TestCheckResourceAttrSet("stackit_kubernetes_cluster.example", "node_pools.0.os_version"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.minimum", "1"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.maximum", "2"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.max_surge", "1"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.max_unavailable", "1"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.volume_type", "storage_premium_perf1"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.volume_size_gb", "20"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.container_runtime", "containerd"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.zones.0", "eu01-m"),
					resource.TestCheckResourceAttrSet("stackit_kubernetes_cluster.example", "status"),
					resource.TestCheckResourceAttrSet("stackit_kubernetes_cluster.example", "kube_config"),
				),
			},
			// check update of multiple configuration options
			{
				Config: configExtended(name, "new-nodepl", "c1.2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "name", name),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "allow_privileged_containers", "false"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.name", "new-nodepl"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.machine_type", "c1.2"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.os_name", "flatcar"),
					resource.TestCheckResourceAttrSet("stackit_kubernetes_cluster.example", "node_pools.0.os_version"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.minimum", "1"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.maximum", "1"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.max_surge", "1"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.max_unavailable", "1"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.volume_type", "storage_premium_perf1"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.volume_size_gb", "20"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.container_runtime", "containerd"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.zones.0", "eu01-1"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.labels.az", "1"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.labels.name", "example-np"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.taints.0.effect", "PreferNoSchedule"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.taints.0.key", "key2"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.taints.0.value", "value1"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "maintenance.enable_kubernetes_version_updates", "true"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "maintenance.enable_machine_image_version_updates", "true"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "hibernations.0.start", "15 6 * * *"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "hibernations.0.end", "30 20 * * *"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "hibernations.0.timezone", "Europe/Berlin"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "extensions.argus.enabled", "false"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "extensions.acl.enabled", "true"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "extensions.acl.allowed_cidrs.0", "185.124.192.0/22"),
					resource.TestCheckResourceAttrSet("stackit_kubernetes_cluster.example", "status"),
					resource.TestCheckResourceAttrSet("stackit_kubernetes_cluster.example", "kube_config"),
				),
			},
			// change machine type
			{
				Config: configExtended(name, "new-nodepl", "c1.3"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "name", name),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.name", "new-nodepl"),
					resource.TestCheckResourceAttr("stackit_kubernetes_cluster.example", "node_pools.0.machine_type", "c1.3"),
				),
			},
			// test import
			{
				ResourceName:            "stackit_kubernetes_cluster.example",
				ImportStateId:           fmt.Sprintf("%s,%s", common.GetAcceptanceTestsProjectID(), name),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status", "extensions", "kubernetes_version"},
			},
		},
	})
}

func configMinimal(name string) string {
	return fmt.Sprintf(`

resource "stackit_kubernetes_project" "example" {
	project_id = "%s"
}

resource "stackit_kubernetes_cluster" "example" {
	kubernetes_project_id         = stackit_kubernetes_project.example.id
	name               			  = "%s"
	
	node_pools = [{
		name         = "example-np"
		machine_type = "c1.2"
	}]
}
	  `,
		common.GetAcceptanceTestsProjectID(),
		name,
	)
}

func configExtended(name, nodepoolName, machineType string) string {
	return fmt.Sprintf(`
resource "stackit_kubernetes_project" "example" {
	project_id = "%s"
}

resource "stackit_kubernetes_cluster" "example" {
	kubernetes_project_id         = stackit_kubernetes_project.example.id
	name               			  = "%s"
	allow_privileged_containers   = false
	
	node_pools = [{
		name         = "%s"
		machine_type = "%s"
		zones        = ["eu01-1"]
		maximum      = 1
	
		labels = {
		  "az"   = "1"
		  "name" = "example-np"
		}
	
		taints = [{
		  effect = "PreferNoSchedule"
		  key    = "key2"
		  value  = "value1"
		}]
	}]

	maintenance = {
		enable_kubernetes_version_updates    = true
		enable_machine_image_version_updates = true
		start                                = "0000-01-01T23:00:00Z"
		end                                  = "0000-01-01T23:30:00Z"
	}

	hibernations = [{
		start    = "15 6 * * *"
		end      = "30 20 * * *"
		timezone = "Europe/Berlin"
	}]

	extensions = {
		argus = {
			enabled = false
		}
		acl = {
			enabled = true
			allowed_cidrs = ["185.124.192.0/22"]
		}
	}
}
	  `,
		common.GetAcceptanceTestsProjectID(),
		name,
		nodepoolName,
		machineType,
	)
}
