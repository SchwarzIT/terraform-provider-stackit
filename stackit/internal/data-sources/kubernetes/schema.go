package kubernetes

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/kubernetes/clusters"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetSchema returns the terraform schema structure
func (r DataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Data source for kubernetes clusters",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Specifies the resource ID",
				Type:        types.StringType,
				Computed:    true,
			},
			"name": {
				Description: "Specifies the cluster name",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.StringWith(clusters.ValidateClusterName, "validate cluster name"),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"project_id": {
				Description: "The project ID the cluster runs in",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.ProjectID(),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"kubernetes_version": {
				Description: "Kubernetes version",
				Type:        types.StringType,
				Computed:    true,
			},
			"kubernetes_version_used": {
				Description: "Full Kubernetes version used. For the data source, it'll always match `kubernetes_version`",
				Type:        types.StringType,
				Computed:    true,
			},
			"allow_privileged_containers": {
				Description: "Are containers allowed to run in privileged mode?",
				Type:        types.BoolType,
				Computed:    true,
			},

			"node_pools": {
				Description: "One or more `node_pool` blocks",
				Computed:    true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Description: "The name of the node pool",
						Type:        types.StringType,
						Computed:    true,
					},
					"machine_type": {
						Description: "The machine type",
						Type:        types.StringType,
						Computed:    true,
					},
					"os_name": {
						Description: "The name of the OS image",
						Type:        types.StringType,
						Computed:    true,
					},
					"os_version": {
						Description: "The OS image version",
						Type:        types.StringType,
						Computed:    true,
					},
					"minimum": {
						Description: "Minimum nodes in the pool",
						Type:        types.Int64Type,
						Computed:    true,
					},
					"maximum": {
						Description: "Maximum nodes in the pool",
						Type:        types.Int64Type,
						Computed:    true,
					},
					"max_surge": {
						Description: "The maximum number of nodes upgraded simultaneously",
						Type:        types.Int64Type,
						Computed:    true,
					},
					"max_unavailable": {
						Description: "The maximum number of nodes unavailable during upgraded",
						Type:        types.Int64Type,
						Computed:    true,
					},
					"volume_type": {
						Description: "Specifies the volume type",
						Type:        types.StringType,
						Computed:    true,
					},
					"volume_size_gb": {
						Description: "The volume size in GB",
						Type:        types.Int64Type,
						Computed:    true,
					},
					"labels": {
						Description: "Labels added to each node",
						Type: types.MapType{
							ElemType: types.StringType,
						},
						Computed: true,
					},
					"taints": {
						Description: "Taint blocks",
						Computed:    true,
						Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
							"effect": {
								Description: "The taint effect",
								Type:        types.StringType,
								Computed:    true,
							},
							"key": {
								Description: "Taint key applied to a node",
								Type:        types.StringType,
								Computed:    true,
							},
							"value": {
								Description: "Taint value corresponding to the taint key",
								Type:        types.StringType,
								Computed:    true,
							},
						}),
					},
					"container_runtime": {
						Description: "The container runtime",
						Type:        types.StringType,
						Computed:    true,
					},
					"zones": {
						Description: "List of availability zones",
						Type:        types.ListType{ElemType: types.StringType},
						Computed:    true,
					},
				}),
			},

			"maintenance": {
				Description: "A single maintenance block as defined below",
				Computed:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"enable_kubernetes_version_updates": {
						Description: "Flag to enable/disable auto-updates of the Kubernetes version",
						Type:        types.BoolType,
						Computed:    true,
					},
					"enable_machine_image_version_updates": {
						Description: "Flag to enable/disable auto-updates of the OS image version",
						Type:        types.BoolType,
						Computed:    true,
					},
					"start": {
						Description: "RFC3339 Date time for maintenance window start. i.e. `2019-08-24T23:00:00Z`",
						Type:        types.StringType,
						Computed:    true,
					},
					"end": {
						Description: "RFC3339 Date time for maintenance window end. i.e. `2019-08-24T23:30:00Z`",
						Type:        types.StringType,
						Computed:    true,
					},
				}),
			},

			"hibernations": {
				Description: "One or more hibernation blocks",
				Computed:    true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"start": {
						Description: "Start time of cluster hibernation",
						Type:        types.StringType,
						Computed:    true,
					},
					"end": {
						Description: "End time of hibernation",
						Type:        types.StringType,
						Computed:    true,
					},
					"timezone": {
						Description: "Timezone",
						Type:        types.StringType,
						Computed:    true,
					},
				}),
			},

			"extensions": {
				Description: "A single extensions block",
				Computed:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"argus": {
						Description: "A single argus block",
						Optional:    true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"enabled": {
								Description: "Is argus extension enabled?",
								Type:        types.BoolType,
								Computed:    true,
							},
							"argus_instance_id": {
								Description: "Instance ID of argus",
								Type:        types.StringType,
								Computed:    true,
							},
						}),
					},
				}),
			},

			"status": {
				Description: "The cluster aggregated status",
				Type:        types.StringType,
				Computed:    true,
			},

			"kube_config": {
				Description: "Kube config file used for connecting to the cluster",
				Type:        types.StringType,
				Sensitive:   true,
				Computed:    true,
			},
		},
	}, nil
}
