package cluster

import (
	"context"
	"fmt"
	"regexp"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/kubernetes/v1.0/cluster"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Cluster is the schema model
type Cluster struct {
	ID                        types.String   `tfsdk:"id"`
	Name                      types.String   `tfsdk:"name"`
	KubernetesProjectID       types.String   `tfsdk:"kubernetes_project_id"`
	ProjectID                 types.String   `tfsdk:"project_id"`
	KubernetesVersion         types.String   `tfsdk:"kubernetes_version"`
	KubernetesVersionUsed     types.String   `tfsdk:"kubernetes_version_used"`
	AllowPrivilegedContainers types.Bool     `tfsdk:"allow_privileged_containers"`
	NodePools                 []NodePool     `tfsdk:"node_pools"`
	Maintenance               *Maintenance   `tfsdk:"maintenance"`
	Hibernations              []Hibernation  `tfsdk:"hibernations"`
	Extensions                *Extensions    `tfsdk:"extensions"`
	Status                    types.String   `tfsdk:"status"`
	KubeConfig                types.String   `tfsdk:"kube_config"`
	Timeouts                  timeouts.Value `tfsdk:"timeouts"`
	NetworkID                 types.String   `tfsdk:"network_id"`
}

type NodePool struct {
	Name             types.String `tfsdk:"name"`
	MachineType      types.String `tfsdk:"machine_type"`
	OSName           types.String `tfsdk:"os_name"`
	OSVersion        types.String `tfsdk:"os_version"`
	Minimum          types.Int64  `tfsdk:"minimum"`
	Maximum          types.Int64  `tfsdk:"maximum"`
	MaxSurge         types.Int64  `tfsdk:"max_surge"`
	MaxUnavailable   types.Int64  `tfsdk:"max_unavailable"`
	VolumeType       types.String `tfsdk:"volume_type"`
	VolumeSizeGB     types.Int64  `tfsdk:"volume_size_gb"`
	Labels           types.Map    `tfsdk:"labels"`
	Taints           []Taint      `tfsdk:"taints"`
	ContainerRuntime types.String `tfsdk:"container_runtime"`
	Zones            types.List   `tfsdk:"zones"`
}

type Taint struct {
	Effect types.String `tfsdk:"effect"`
	Key    types.String `tfsdk:"key"`
	Value  types.String `tfsdk:"value"`
}

type Maintenance struct {
	EnableKubernetesVersionUpdates   types.Bool   `tfsdk:"enable_kubernetes_version_updates"`
	EnableMachineImageVersionUpdates types.Bool   `tfsdk:"enable_machine_image_version_updates"`
	Start                            types.String `tfsdk:"start"`
	End                              types.String `tfsdk:"end"`
}

type Hibernation struct {
	Start    types.String `tfsdk:"start"`
	End      types.String `tfsdk:"end"`
	Timezone types.String `tfsdk:"timezone"`
}

type Extensions struct {
	Argus *ArgusExtension `tfsdk:"argus"`
	ACL   *ACL            `tfsdk:"acl"`
}

type ACL struct {
	Enabled      types.Bool `tfsdk:"enabled"`
	AllowedCIDRs types.List `tfsdk:"allowed_cidrs"`
}

type ArgusExtension struct {
	Enabled         types.Bool   `tfsdk:"enabled"`
	ArgusInstanceID types.String `tfsdk:"argus_instance_id"`
}

// Schema returns the terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Manages kubernetes clusters\n%s",
			common.EnvironmentInfo(r.urls),
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
					validate.StringWith(cluster.ValidateClusterName, "validate cluster name"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// TODO: remove in next releases
			"kubernetes_project_id": schema.StringAttribute{
				Description:        "The ID of a `stackit_kubernetes_project` resource",
				DeprecationMessage: "This attribute is deprecated and will be removed in a future version. Please use the `project_id` attribute instead.",
				Optional:           true,
				Computed:           true,
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
				Description: "Kubernetes version. Allowed Options are: `1.25`, `1.26`, or a full version including patch (not recommended).",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					validate.StringWith(clientValidate.SemVer, "validate container runtime"),
				},
				Default: stringdefault.StaticString(DefaultVersion),
			},
			"kubernetes_version_used": schema.StringAttribute{
				Description: "Full Kubernetes version used. For example, if `1.22` was selected, this value may result to `1.22.15`",
				Computed:    true,
			},
			"allow_privileged_containers": schema.BoolAttribute{
				Description:        "Should containers be allowed to run in privileged mode? Default is `true`",
				DeprecationMessage: "This attribute is deprecated starting from v1.25",
				Optional:           true,
				Computed:           true,
				Default:            booldefault.StaticBool(DefaultAllowPrivileged),
			},

			"node_pools": schema.ListNestedAttribute{
				Description: "One or more `node_pool` block as defined below",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Specifies the name of the node pool",
							Required:    true,
							Validators: []validator.String{
								validate.StringWith(cluster.ValidateNodePoolName, "validate node pool name"),
							},
						},
						"machine_type": schema.StringAttribute{
							Description: "The machine type. Accepted options are: `c1.2`, `c1.3`, `c1.4`, `c1.5`, `g1.2`, `g1.3`, `g1.4`, `g1.5`, `m1.2`, `m1.3`, `m1.4`",
							Required:    true,
						},
						"os_name": schema.StringAttribute{
							Description: "The name of the OS image. Only `flatcar` is supported",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(DefaultOSName),
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
							Default:     int64default.StaticInt64(DefaultNodepoolMin),
						},

						"maximum": schema.Int64Attribute{
							Description: "Maximum nodes in the pool. Defaults to 2. (Value must be between 1-100)",
							Optional:    true,
							Computed:    true,
							Default:     int64default.StaticInt64(DefaultNodepoolMax),
						},

						"max_surge": schema.Int64Attribute{
							Description: "The maximum number of nodes upgraded simultaneously. Defaults to 1. (Value must be between 1-10)",
							Optional:    true,
							Computed:    true,
							Default:     int64default.StaticInt64(DefaultNodepoolMaxSurge),
						},
						"max_unavailable": schema.Int64Attribute{
							Description: "The maximum number of nodes unavailable during upgraded. Defaults to 0",
							Optional:    true,
							Computed:    true,
							Default:     int64default.StaticInt64(DefaultNodepoolMaxUnavailable),
						},
						"volume_type": schema.StringAttribute{
							Description: "Specifies the volume type. Defaults to `storage_premium_perf1`. Available options are `storage_premium_perf0`, `storage_premium_perf1`, `storage_premium_perf2`, `storage_premium_perf4`, `storage_premium_perf6`",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(DefaultVolumeType),
						},
						"volume_size_gb": schema.Int64Attribute{
							Description: "The volume size in GB. Default is set to `20`",
							Optional:    true,
							Computed:    true,
							Default:     int64default.StaticInt64(DefaultVolumeSizeGB),
						},
						"labels": schema.MapAttribute{
							Description: "Labels to add to each node",
							ElementType: types.StringType,
							Optional:    true,
							Computed:    true,
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
							Default: stringdefault.StaticString(DefaultCRI),
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
						Description: "RFC3339 Date time for maintenance window start. i.e. `0000-01-01T23:00:00Z`",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(`^(0000-01-01T\d{2}:\d{2}:\d{2}Z)$`), "validate RFC3339 date time that starts with 0000-01-01"),
						},
					},
					"end": schema.StringAttribute{
						Description: "RFC3339 Date time for maintenance window end. i.e. `0000-01-01T23:30:00Z`",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(`^(0000-01-01T\d{2}:\d{2}:\d{2}Z)$`), "validate RFC3339 date time that starts with 0000-01-01"),
						},
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
					"acl": schema.SingleNestedAttribute{
						Description: "Cluster access control configuration",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Is ACL enabled? Defaults to `false`",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"allowed_cidrs": schema.ListAttribute{
								Description: "Specify a list of CIDRs to whitelist",
								ElementType: types.StringType,
								Optional:    true,
								Computed:    true,
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

			"timeouts": common.Timeouts(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),

			"network_id": schema.StringAttribute{
				Description: "Specifies the ID of the Network the SKE-Nodes should be created in",
				Required:    false,
				Computed:    false,
				Optional:    true,
				Validators: []validator.String{
					validate.NetworkID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}
