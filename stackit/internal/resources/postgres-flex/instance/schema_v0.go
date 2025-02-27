package postgresinstance

import (
	"context"
	"fmt"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getSchemaV0(ctx context.Context) *schema.Schema {
	return &schema.Schema{
		Version: 1,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the resource ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Specifies the instance name.",
				Required:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID the instance runs in. Changing this value requires the resource to be recreated.",
				Required:    true,
			},
			"machine_type": schema.StringAttribute{
				Description: "The Machine Type. Available options: `2.4` (2 CPU, 4 Memory), `2.16` (2 CPU, 16 Memory), `4.8` (4 CPU, 8 Memory), `4.32` (4 CPU, 32 Memory), `8.16` (8 CPU, 16 Memory), `16.32` (16 CPU, 32 Memory), `16.128` (16 CPU, 128 Memory)",
				Required:    true,
			},
			"version": schema.StringAttribute{
				Description: "Postgres version. Options: `12`, `13`, `14`. Changing this value requires the resource to be recreated.",
				Optional:    true,
				Computed:    true,
			},
			"replicas": schema.Int64Attribute{
				Description: "Number of replicas (Default is `1`).",
				Optional:    true,
				Computed:    true,
			},
			"backup_schedule": schema.StringAttribute{
				Description: "Specifies the backup schedule (cron style)",
				Optional:    true,
				Computed:    true,
			},
			"storage": schema.SingleNestedAttribute{
				Description: "A single `storage` block as defined below.",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"class": schema.StringAttribute{
						Description: "Specifies the storage class. Available option: `premium-perf6-stackit`",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(DefaultStorageClass),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"size": schema.Int64Attribute{
						Description: "The storage size in GB (min of 5 is required)",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(DefaultStorageSize),
					},
				},
			},
			"options": schema.MapAttribute{
				Description: "Specifies postgres instance options",
				ElementType: types.StringType,
				Optional:    true,
			},
			"labels": schema.MapAttribute{
				Description: "Instance Labels",
				ElementType: types.StringType,
				Optional:    true,
			},
			"acl": schema.ListAttribute{
				Description: fmt.Sprintf("Whitelist IP address ranges. Default is %v", common.KnownRanges),
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"timeouts": common.Timeouts(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

func upgradeV0(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	type InstanceV0 struct {
		ID             types.String      `tfsdk:"id"`
		Name           types.String      `tfsdk:"name"`
		ProjectID      types.String      `tfsdk:"project_id"`
		MachineType    types.String      `tfsdk:"machine_type"`
		Version        types.String      `tfsdk:"version"`
		Replicas       types.Int64       `tfsdk:"replicas"`
		BackupSchedule types.String      `tfsdk:"backup_schedule"`
		Options        map[string]string `tfsdk:"options"`
		Labels         map[string]string `tfsdk:"labels"`
		ACL            types.List        `tfsdk:"acl"`
		Storage        types.Object      `tfsdk:"storage"`
		Timeouts       timeouts.Value    `tfsdk:"timeouts"`
	}

	var oldState InstanceV0

	diags := req.State.Get(ctx, &oldState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	acl, ds := types.SetValueFrom(ctx, types.StringType, oldState.ACL.Elements())
	if ds.HasError() {
		resp.Diagnostics.Append(ds...)
		return
	}

	newState := Instance{
		ID:             oldState.ID,
		Name:           oldState.Name,
		ProjectID:      oldState.ProjectID,
		MachineType:    oldState.MachineType,
		Version:        oldState.Version,
		Replicas:       oldState.Replicas,
		BackupSchedule: oldState.BackupSchedule,
		Options:        oldState.Options,
		Labels:         oldState.Labels,
		ACL:            acl,
		Storage:        oldState.Storage,
		Timeouts:       oldState.Timeouts,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}
