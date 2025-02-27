package instance

import (
	"context"
	"fmt"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getSchemaV0(ctx context.Context) *schema.Schema {
	return &schema.Schema{
		Version: 1,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the resource ID.",
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
				Description: "The Machine Type. Available options: `1.1` (1 CPU, 1 Memory), `1.2` (1 CPU, 2 Memory), `1.4` (1 CPU, 4 Memory), `1.8` (1 CPU, 8 Memory), `2.4` (2 CPU, 4 Memory), `2.8` (2 CPU, 8 Memory), `2.16` (2 CPU, 16 Memory), `4.8` (4 CPU, 8 Memory), `4.16` (4 CPU, 16 Memory), `4.32` (4 CPU, 32 Memory), `8.16` (8 CPU, 16 Memory), `8.32` (8 CPU, 32 Memory), `8.64` (8 CPU, 64 Memory), `16.32` (16 CPU, 32 Memory), `16.64` (16 CPU, 64 Memory)",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The service type. Available options: `Single`, `Replica`, `Sharded`. Changing this value requires the resource to be recreated.",
				Optional:    true,
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "MongoDB version. Version `5.0`, `6.0`, `7.0` are supported. ",
				Optional:    true,
				Computed:    true,
			},
			"replicas": schema.Int64Attribute{
				Description: fmt.Sprintf("Number of replicas (Default is `%d`).", DefaultReplicas),
				Optional:    true,
				Computed:    true,
			},
			"backup_schedule": schema.StringAttribute{
				Description: "Specifies the backup schedule (cron style).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(DefaultBackupSchedule),
			},
			"storage": schema.SingleNestedAttribute{
				Description: "A single `storage` block as defined below.",
				Optional:    true,
				Computed:    true,
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
		Type           types.String      `tfsdk:"type"`
		MachineType    types.String      `tfsdk:"machine_type"` // aka FlavorID
		Version        types.String      `tfsdk:"version"`
		Replicas       types.Int64       `tfsdk:"replicas"`
		BackupSchedule types.String      `tfsdk:"backup_schedule"`
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
		Type:           oldState.Type,
		MachineType:    oldState.MachineType,
		Version:        oldState.Version,
		Replicas:       oldState.Replicas,
		BackupSchedule: oldState.BackupSchedule,
		Labels:         oldState.Labels,
		ACL:            acl,
		Storage:        oldState.Storage,
		Timeouts:       oldState.Timeouts,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}
