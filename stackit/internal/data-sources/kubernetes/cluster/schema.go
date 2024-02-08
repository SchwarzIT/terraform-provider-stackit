package cluster

import (
	"context"
	"fmt"

	clientCluster "github.com/SchwarzIT/community-stackit-go-client/pkg/services/kubernetes/v1.0/cluster"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/resources/kubernetes/cluster"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Cluster struct {
	ID                    types.String          `tfsdk:"id"`
	Name                  types.String          `tfsdk:"name"`
	ProjectID             types.String          `tfsdk:"project_id"`
	KubernetesProjectID   types.String          `tfsdk:"kubernetes_project_id"`
	KubernetesVersion     types.String          `tfsdk:"kubernetes_version"`
	KubernetesVersionUsed types.String          `tfsdk:"kubernetes_version_used"`
	NodePools             []cluster.NodePool    `tfsdk:"node_pools"`
	Maintenance           *cluster.Maintenance  `tfsdk:"maintenance"`
	Hibernations          []cluster.Hibernation `tfsdk:"hibernations"`
	Extensions            *cluster.Extensions   `tfsdk:"extensions"`
	Status                types.String          `tfsdk:"status"`
	KubeConfig            types.String          `tfsdk:"kube_config"`
	NetworkID             types.String          `tfsdk:"network_id"`
}

// Schema returns the terraform schema structure
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Data source for STACKIT Kubernetes Engine (SKE) clusters\n%s",
			common.EnvironmentInfo(d.urls),
		),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the resource ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Specifies the cluster name (lower case, alphanumeric, hypens allowed, up to 11 chars)",
				Required:    true,
				Validators: []validator.String{
					validate.StringWith(clientCluster.ValidateClusterName, "validate cluster name"),
				},
			},
			// TODO: remove in next releases
			"kubernetes_project_id": schema.StringAttribute{
				Description:        "The ID of a `stackit_kubernetes_project` resource",
				DeprecationMessage: "This attribute is being deprecated. Use `project_id` instead",
				Optional:           true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "The project UUID.",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
			},
			"kubernetes_version": schema.StringAttribute{
				Description: "Kubernetes version. ",
				Computed:    true,
			},
			"kubernetes_version_used": schema.StringAttribute{
				Description: "Full Kubernetes version used. For example, if `1.22` was selected, this value may result to `1.22.15`",
				Computed:    true,
			},
			"node_pools": schema.ListNestedAttribute{
				Description: "One or more `node_pool` block as defined below",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Specifies the name of the node pool",
							Required:    true,
						},
						"machine_type": schema.StringAttribute{
							Description: "The machine type.",
							Required:    true,
						},
						"os_name": schema.StringAttribute{
							Description: "The name of the OS image.",
							Computed:    true,
						},
						"os_version": schema.StringAttribute{
							Description: "The OS image version.",
							Computed:    true,
						},
						"minimum": schema.Int64Attribute{
							Description: "Minimum nodes in the pool.",
							Computed:    true,
						},

						"maximum": schema.Int64Attribute{
							Description: "Maximum nodes in the pool.",
							Computed:    true,
						},

						"max_surge": schema.Int64Attribute{
							Description: "The maximum number of nodes upgraded simultaneously.",
							Computed:    true,
						},
						"max_unavailable": schema.Int64Attribute{
							Description: "The maximum number of nodes unavailable during upgraded.",
							Computed:    true,
						},
						"volume_type": schema.StringAttribute{
							Description: "Specifies the volume type.",
							Computed:    true,
						},
						"volume_size_gb": schema.Int64Attribute{
							Description: "The volume size in GB.",
							Computed:    true,
						},
						"labels": schema.MapAttribute{
							Description: "Labels to add to each node",
							ElementType: types.StringType,
							Computed:    true,
						},
						"taints": schema.ListNestedAttribute{
							Description: "Specifies a taint list as defined below",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"effect": schema.StringAttribute{
										Description: "The taint effect. Only `PreferNoSchedule` is supported at the moment",
										Computed:    true,
									},
									"key": schema.StringAttribute{
										Description: "Taint key to be applied to a node",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "Taint value corresponding to the taint key",
										Computed:    true,
									},
								},
							},
						},
						"container_runtime": schema.StringAttribute{
							Description: "Specifies the container runtime.",
							Computed:    true,
						},
						"zones": schema.ListAttribute{
							Description: "Specify a list of availability zones.",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
			"maintenance": schema.SingleNestedAttribute{
				Description: "A single maintenance block as defined below",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"enable_kubernetes_version_updates": schema.BoolAttribute{
						Description: "Flag to enable/disable auto-updates of the Kubernetes version",
						Computed:    true,
					},
					"enable_machine_image_version_updates": schema.BoolAttribute{
						Description: "Flag to enable/disable auto-updates of the OS image version",
						Computed:    true,
					},
					"start": schema.StringAttribute{
						Description: "RFC3339 Date time for maintenance window start. i.e. `2019-08-24T23:00:00Z`",
						Computed:    true,
					},
					"end": schema.StringAttribute{
						Description: "RFC3339 Date time for maintenance window end. i.e. `2019-08-24T23:30:00Z`",
						Computed:    true,
					},
				},
			},

			"hibernations": schema.ListNestedAttribute{
				Description: "One or more hibernation block as defined below",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"start": schema.StringAttribute{
							Description: "Start time of cluster hibernation, in crontab syntax, i.e. `0 18 * * *` for starting everyday at 6pm",
							Computed:    true,
						},
						"end": schema.StringAttribute{
							Description: "End time of hibernation, in crontab syntax, i.e. `0 8 * * *` for waking up the cluster at 8am",
							Computed:    true,
						},
						"timezone": schema.StringAttribute{
							Description: "Timezone name corresponding to a file in the IANA Time Zone database. i.e. `Europe/Berlin`",
							Computed:    true,
						},
					},
				},
			},

			"extensions": schema.SingleNestedAttribute{
				Description: "A single extensions block as defined below",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"argus": schema.SingleNestedAttribute{
						Description: "A single argus block as defined below",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Flag to enable/disable argus extensions.",
								Computed:    true,
							},
							"argus_instance_id": schema.StringAttribute{
								Description: "Instance ID of argus",
								Computed:    true,
							},
						},
					},
					"acl": schema.SingleNestedAttribute{
						Description: "Cluster access control configuration",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Is ACL enabled? Defaults to `false`",
								Computed:    true,
							},
							"allowed_cidrs": schema.ListAttribute{
								Description: "Specify a list of CIDRs to whitelist",
								ElementType: types.StringType,
								Computed:    true,
							},
						},
					},
				},
			},

			"status": schema.StringAttribute{
				Description: "The cluster's aggregated status",
				Computed:    true,
			},

			"kube_config": schema.StringAttribute{
				Description: "Kube config file used for connecting to the cluster",
				Sensitive:   true,
				Computed:    true,
			},

			"network_id": schema.StringAttribute{
				Description: "Specifies the ID of the SNA-Network created",
				Required:    false,
				Computed:    false,
				Optional:    true,
			},
		},
	}
}
