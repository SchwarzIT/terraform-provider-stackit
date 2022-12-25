package cluster

import (
	"context"

	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/kubernetes/v1.0/generated/cluster"
	clientValidate "github.com/SchwarzIT/community-stackit-go-client/pkg/validate"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/modifiers"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Cluster is the schema model
type Cluster struct {
	ID                        types.String  `tfsdk:"id"`
	Name                      types.String  `tfsdk:"name"`
	KubernetesProjectID       types.String  `tfsdk:"kubernetes_project_id"`
	KubernetesVersion         types.String  `tfsdk:"kubernetes_version"`
	KubernetesVersionUsed     types.String  `tfsdk:"kubernetes_version_used"`
	AllowPrivilegedContainers types.Bool    `tfsdk:"allow_privileged_containers"`
	NodePools                 []NodePool    `tfsdk:"node_pools"`
	Maintenance               *Maintenance  `tfsdk:"maintenance"`
	Hibernations              []Hibernation `tfsdk:"hibernations"`
	Extensions                *Extensions   `tfsdk:"extensions"`
	Status                    types.String  `tfsdk:"status"`
	KubeConfig                types.String  `tfsdk:"kube_config"`
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
}

type ArgusExtension struct {
	Enabled         types.Bool   `tfsdk:"enabled"`
	ArgusInstanceID types.String `tfsdk:"argus_instance_id"`
}

// GetSchema returns the terraform schema structure
func (r *Resource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages kubernetes clusters",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Specifies the resource ID",
				Type:        types.StringType,
				Computed:    true,
			},
			"name": {
				Description: "Specifies the cluster name (lower case, alphanumeric, hypens allowed, up to 11 chars)",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.StringWith(cluster.ValidateClusterName, "validate cluster name"),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"kubernetes_project_id": {
				Description: "The ID of a `stackit_kubernetes_project` resource",
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
				Description: "Kubernetes version. Allowed Options are: `1.22`, `1.23`, `1.24`, or a full version including patch (not recommended).",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
				Validators: []tfsdk.AttributeValidator{
					validate.StringWith(clientValidate.SemVer, "validate container runtime"),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.StringDefault(default_version),
				},
			},
			"kubernetes_version_used": {
				Description: "Full Kubernetes version used. For example, if `1.22` was selected, this value may result to `1.22.15`",
				Type:        types.StringType,
				Computed:    true,
			},
			"allow_privileged_containers": {
				Description: "Should containers be allowed to run in privileged mode? Default is `true`",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.BoolDefault(default_allow_privileged),
				},
			},

			"node_pools": {
				Description: "One or more `node_pool` block as defined below",
				Optional:    true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Description: "Specifies the name of the node pool",
						Type:        types.StringType,
						Required:    true,
					},
					"machine_type": {
						Description: "The machine type. Accepted options are: `c1.2`, `c1.3`, `c1.4`, `c1.5`, `g1.2`, `g1.3`, `g1.4`, `g1.5`, `m1.2`, `m1.3`, `m1.4`",
						Type:        types.StringType,
						Required:    true,
					},
					"os_name": {
						Description: "The name of the OS image. Only `flatcar` is supported",
						Type:        types.StringType,
						Optional:    true,
						Computed:    true,
						PlanModifiers: []tfsdk.AttributePlanModifier{
							modifiers.StringDefault(default_os_name),
						},
					},
					"os_version": {
						Description: "The OS image version.",
						Type:        types.StringType,
						Optional:    true,
						Computed:    true,
					},
					"minimum": {
						Description: "Minimum nodes in the pool. Defaults to 1. (Value must be between 1-100)",
						Type:        types.Int64Type,
						Optional:    true,
						Computed:    true,
						PlanModifiers: []tfsdk.AttributePlanModifier{
							modifiers.Int64Default(default_nodepool_min),
						},
					},
					"maximum": {
						Description: "Maximum nodes in the pool. Defaults to 2. (Value must be between 1-100)",
						Type:        types.Int64Type,
						Optional:    true,
						Computed:    true,
						PlanModifiers: []tfsdk.AttributePlanModifier{
							modifiers.Int64Default(default_nodepool_max),
						},
					},
					"max_surge": {
						Description: "The maximum number of nodes upgraded simultaneously. Defaults to 1. (Value must be between 1-10)",
						Type:        types.Int64Type,
						Optional:    true,
						Computed:    true,
						PlanModifiers: []tfsdk.AttributePlanModifier{
							modifiers.Int64Default(default_nodepool_max_surge),
						},
					},
					"max_unavailable": {
						Description: "The maximum number of nodes unavailable during upgraded. Defaults to 1",
						Type:        types.Int64Type,
						Optional:    true,
						Computed:    true,
						PlanModifiers: []tfsdk.AttributePlanModifier{
							modifiers.Int64Default(default_nodepool_max_unavailable),
						},
					},
					"volume_type": {
						Description: "Specifies the volume type. Defaults to `storage_premium_perf1`. Available options are `storage_premium_perf0`, `storage_premium_perf1`, `storage_premium_perf2`, `storage_premium_perf4`, `storage_premium_perf6`",
						Type:        types.StringType,
						Optional:    true,
						Computed:    true,
						PlanModifiers: []tfsdk.AttributePlanModifier{
							modifiers.StringDefault(default_volume_type),
						},
					},
					"volume_size_gb": {
						Description: "The volume size in GB. Default is set to `20`",
						Type:        types.Int64Type,
						Optional:    true,
						Computed:    true,
						PlanModifiers: []tfsdk.AttributePlanModifier{
							modifiers.Int64Default(default_volume_size_gb),
						},
					},
					"labels": {
						Description: "Labels to add to each node",
						Type: types.MapType{
							ElemType: types.StringType,
						},
						Optional: true,
					},
					"taints": {
						Description: "Specifies a taint list as defined below",
						Optional:    true,
						Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
							"effect": {
								Description: "The taint effect. Only `PreferNoSchedule` is supported at the moment",
								Type:        types.StringType,
								Required:    true,
							},
							"key": {
								Description: "Taint key to be applied to a node",
								Type:        types.StringType,
								Required:    true,
							},
							"value": {
								Description: "Taint value corresponding to the taint key",
								Type:        types.StringType,
								Optional:    true,
							},
						}),
					},
					"container_runtime": {
						Description: "Specifies the container runtime. Defaults to `containerd`. Allowed options are `docker`, `containerd`",
						Type:        types.StringType,
						Optional:    true,
						Computed:    true,
						Validators: []tfsdk.AttributeValidator{
							validate.StringWith(func(v string) error {
								n := cluster.CRIName(v)
								cri := cluster.CRI{Name: &n}
								return cluster.ValidateCRI(&cri)
							}, "validate container runtime"),
						},
						PlanModifiers: []tfsdk.AttributePlanModifier{
							modifiers.StringDefault(default_cri),
						},
					},
					"zones": {
						Description: "Specify a list of availability zones. Accepted options are `eu01-m` for metro, or `eu01-1`, `eu01-2`, `eu01-3`",
						Type:        types.ListType{ElemType: types.StringType},
						Optional:    true,
						Computed:    true,
					},
				}),
			},

			"maintenance": {
				Description: "A single maintenance block as defined below",
				Optional:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"enable_kubernetes_version_updates": {
						Description: "Flag to enable/disable auto-updates of the Kubernetes version",
						Type:        types.BoolType,
						Required:    true,
					},
					"enable_machine_image_version_updates": {
						Description: "Flag to enable/disable auto-updates of the OS image version",
						Type:        types.BoolType,
						Required:    true,
					},
					"start": {
						Description: "RFC3339 Date time for maintenance window start. i.e. `2019-08-24T23:00:00Z`",
						Type:        types.StringType,
						Required:    true,
					},
					"end": {
						Description: "RFC3339 Date time for maintenance window end. i.e. `2019-08-24T23:30:00Z`",
						Type:        types.StringType,
						Required:    true,
					},
				}),
			},

			"hibernations": {
				Description: "One or more hibernation block as defined below",
				Optional:    true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"start": {
						Description: "Start time of cluster hibernation, in crontab syntax, i.e. `0 18 * * *` for starting everyday at 6pm",
						Type:        types.StringType,
						Required:    true,
					},
					"end": {
						Description: "End time of hibernation, in crontab syntax, i.e. `0 8 * * *` for waking up the cluster at 8am",
						Type:        types.StringType,
						Required:    true,
					},
					"timezone": {
						Description: "Timezone name corresponding to a file in the IANA Time Zone database. i.e. `Europe/Berlin`",
						Type:        types.StringType,
						Optional:    true,
					},
				}),
			},

			"extensions": {
				Description: "A single extensions block as defined below",
				Optional:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"argus": {
						Description: "A single argus block as defined below",
						Optional:    true,
						Computed:    true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"enabled": {
								Description: "Flag to enable/disable argus extensions. Defaults to `false`",
								Type:        types.BoolType,
								Optional:    true,
								Computed:    true,
							},
							"argus_instance_id": {
								Description: "Instance ID of argus, Required when enabled is set to `true`",
								Type:        types.StringType,
								Optional:    true,
							},
						}),
					},
				}),
			},

			"status": {
				Description: "The cluster's aggregated status",
				Type:        types.StringType,
				Computed:    true,
				Required:    false,
				Optional:    false,
			},

			"kube_config": {
				Description: "Kube config file used for connecting to the cluster",
				Type:        types.StringType,
				Sensitive:   true,
				Computed:    true,
				Required:    false,
				Optional:    false,
			},
		},
	}, nil
}
