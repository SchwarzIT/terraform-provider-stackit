package cluster

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/kubernetes/v1.0/generated/cluster"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/modifiers"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Cluster is the schema model
type Cluster struct {
	ID                        types.String  `schema:"id"`
	Name                      types.String  `schema:"name"`
	KubernetesProjectID       types.String  `schema:"kubernetes_project_id"`
	KubernetesVersion         types.String  `schema:"kubernetes_version"`
	KubernetesVersionUsed     types.String  `schema:"kubernetes_version_used"`
	AllowPrivilegedContainers types.Bool    `schema:"allow_privileged_containers"`
	NodePools                 []NodePool    `schema:"node_pools"`
	Maintenance               *Maintenance  `schema:"maintenance"`
	Hibernations              []Hibernation `schema:"hibernations"`
	Extensions                *Extensions   `schema:"extensions"`
	Status                    types.String  `schema:"status"`
	KubeConfig                types.String  `schema:"kube_config"`
}

type NodePool struct {
	Name             types.String `schema:"name"`
	MachineType      types.String `schema:"machine_type"`
	OSName           types.String `schema:"os_name"`
	OSVersion        types.String `schema:"os_version"`
	Minimum          types.Int64  `schema:"minimum"`
	Maximum          types.Int64  `schema:"maximum"`
	MaxSurge         types.Int64  `schema:"max_surge"`
	MaxUnavailable   types.Int64  `schema:"max_unavailable"`
	VolumeType       types.String `schema:"volume_type"`
	VolumeSizeGB     types.Int64  `schema:"volume_size_gb"`
	Labels           types.Map    `schema:"labels"`
	Taints           []Taint      `schema:"taints"`
	ContainerRuntime types.String `schema:"container_runtime"`
	Zones            types.List   `schema:"zones"`
}

type Taint struct {
	Effect types.String `schema:"effect"`
	Key    types.String `schema:"key"`
	Value  types.String `schema:"value"`
}

type Maintenance struct {
	EnableKubernetesVersionUpdates   types.Bool   `schema:"enable_kubernetes_version_updates"`
	EnableMachineImageVersionUpdates types.Bool   `schema:"enable_machine_image_version_updates"`
	Start                            types.String `schema:"start"`
	End                              types.String `schema:"end"`
}

type Hibernation struct {
	Start    types.String `schema:"start"`
	End      types.String `schema:"end"`
	Timezone types.String `schema:"timezone"`
}

type Extensions struct {
	Argus *ArgusExtension `schema:"argus"`
}

type ArgusExtension struct {
	Enabled         types.Bool   `schema:"enabled"`
	ArgusInstanceID types.String `schema:"argus_instance_id"`
}

// GetSchema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages kubernetes clusters",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the resource ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Specifies the cluster name (lower case, alphanumeric, hypens allowed, up to 11 chars)",
				Required:    true,
				Validators: []validator.String{
					validate.StringWith(cluster.ValidateClusterName, "validate cluster name"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"kubernetes_project_id": schema.StringAttribute{
				Description: "The ID of a `stackit_kubernetes_project` resource",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"kubernetes_version": schema.StringAttribute{
				Description: "Kubernetes version. Allowed Options are: `1.22`, `1.23`, `1.24`, or a full version including patch (not recommended).",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					validate.StringWith(clientValidate.SemVer, "validate container runtime"),
				},
				PlanModifiers: []planmodifier.String{
					modifiers.StringDefault(default_version),
				},
			},
			"kubernetes_version_used": schema.StringAttribute{
				Description: "Full Kubernetes version used. For example, if `1.22` was selected, this value may result to `1.22.15`",
				Computed:    true,
			},
			"allow_privileged_containers": schema.BoolAttribute{
				Description: "Should containers be allowed to run in privileged mode? Default is `true`",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					modifiers.BoolDefault(default_allow_privileged),
				},
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
							Description: "The machine type. Accepted options are: `c1.2`, `c1.3`, `c1.4`, `c1.5`, `g1.2`, `g1.3`, `g1.4`, `g1.5`, `m1.2`, `m1.3`, `m1.4`",
							Required:    true,
						},
						"os_name": schema.StringAttribute{
							Description: "The name of the OS image. Only `flatcar` is supported",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								modifiers.StringDefault(default_os_name),
							},
						},
						"os_version": schema.StringAttribute{
							Description: "The OS image version.",
							Optional:    true,
							Computed:    true,
						},
						"minimum": schema.Int64Attribute{
							Description: "Minimum nodes in the pool. Defaults to 1. (Value must be between 1-100)",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Int64{
								modifiers.Int64Default(default_nodepool_min),
							},
						},

						"maximum": schema.Int64Attribute{
							Description: "Maximum nodes in the pool. Defaults to 2. (Value must be between 1-100)",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Int64{
								modifiers.Int64Default(default_nodepool_max),
							},
						},

						"max_surge": schema.Int64Attribute{
							Description: "The maximum number of nodes upgraded simultaneously. Defaults to 1. (Value must be between 1-10)",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Int64{
								modifiers.Int64Default(default_nodepool_max_surge),
							},
						},
						"max_unavailable": schema.Int64Attribute{
							Description: "The maximum number of nodes unavailable during upgraded. Defaults to 1",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Int64{
								modifiers.Int64Default(default_nodepool_max_unavailable),
							},
						},
						"volume_type": schema.StringAttribute{
							Description: "Specifies the volume type. Defaults to `storage_premium_perf1`. Available options are `storage_premium_perf0`, `storage_premium_perf1`, `storage_premium_perf2`, `storage_premium_perf4`, `storage_premium_perf6`",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								modifiers.StringDefault(default_volume_type),
							},
						},
						"volume_size_gb": schema.Int64Attribute{
							Description: "The volume size in GB. Default is set to `20`",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Int64{
								modifiers.Int64Default(default_volume_size_gb),
							},
						},
						"labels": schema.MapAttribute{
							Description: "Labels to add to each node",
							ElementType: types.StringType,
							Optional:    true,
						},
						"taints": schema.ListNestedAttribute{
							Description: "Specifies a taint list as defined below",
							Optional:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"effect": schema.StringAttribute{
										Description: "The taint effect. Only `PreferNoSchedule` is supported at the moment",
										Required:    true,
									},
									"key": schema.StringAttribute{
										Description: "Taint key to be applied to a node",
										Required:    true,
									},
									"value": schema.StringAttribute{
										Description: "Taint value corresponding to the taint key",
										Optional:    true,
									},
								},
							},
						},
						"container_runtime": schema.StringAttribute{
							Description: "Specifies the container runtime. Defaults to `containerd`. Allowed options are `docker`, `containerd`",
							Optional:    true,
							Computed:    true,
							Validators: []validator.String{
								validate.StringWith(func(v string) error {
									n := cluster.CRIName(v)
									cri := cluster.CRI{Name: &n}
									return cluster.ValidateCRI(&cri)
								}, "validate container runtime"),
							},
							PlanModifiers: []planmodifier.String{
								modifiers.StringDefault(default_cri),
							},
						},
						"zones": schema.ListAttribute{
							Description: "Specify a list of availability zones. Accepted options are `eu01-m` for metro, or `eu01-1`, `eu01-2`, `eu01-3`",
							ElementType: types.StringType,
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"maintenance": schema.SingleNestedAttribute{
				Description: "A single maintenance block as defined below",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"enable_kubernetes_version_updates": schema.BoolAttribute{
						Description: "Flag to enable/disable auto-updates of the Kubernetes version",
						Required:    true,
					},
					"enable_machine_image_version_updates": schema.BoolAttribute{
						Description: "Flag to enable/disable auto-updates of the OS image version",
						Required:    true,
					},
					"start": schema.StringAttribute{
						Description: "RFC3339 Date time for maintenance window start. i.e. `2019-08-24T23:00:00Z`",
						Required:    true,
					},
					"end": schema.StringAttribute{
						Description: "RFC3339 Date time for maintenance window end. i.e. `2019-08-24T23:30:00Z`",
						Required:    true,
					},
				},
			},

			"hibernations": schema.ListNestedAttribute{
				Description: "One or more hibernation block as defined below",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"start": schema.StringAttribute{
							Description: "Start time of cluster hibernation, in crontab syntax, i.e. `0 18 * * *` for starting everyday at 6pm",
							Required:    true,
						},
						"end": schema.StringAttribute{
							Description: "End time of hibernation, in crontab syntax, i.e. `0 8 * * *` for waking up the cluster at 8am",
							Required:    true,
						},
						"timezone": schema.StringAttribute{
							Description: "Timezone name corresponding to a file in the IANA Time Zone database. i.e. `Europe/Berlin`",
							Optional:    true,
						},
					},
				},
			},

			"extensions": schema.SingleNestedAttribute{
				Description: "A single extensions block as defined below",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"argus": schema.SingleNestedAttribute{
						Description: "A single argus block as defined below",
						Optional:    true,
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Flag to enable/disable argus extensions. Defaults to `false`",
								Optional:    true,
								Computed:    true,
							},
							"argus_instance_id": schema.StringAttribute{
								Description: "Instance ID of argus, Required when enabled is set to `true`",
								Optional:    true,
							},
						},
					},
				},
			},

			"status": schema.StringAttribute{
				Description: "The cluster's aggregated status",
				Computed:    true,
				Required:    false,
				Optional:    false,
			},

			"kube_config": schema.StringAttribute{
				Description: "Kube config file used for connecting to the cluster",
				Sensitive:   true,
				Computed:    true,
				Required:    false,
				Optional:    false,
			},
		},
	}
}
