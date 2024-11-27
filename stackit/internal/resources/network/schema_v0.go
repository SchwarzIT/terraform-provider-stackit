package network

import (
	"context"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
				Description: "the name of the network",
				Required:    true,
			},
			"nameservers": schema.ListAttribute{
				Description: "List of DNS Servers/Nameservers.",
				Required:    true,
				ElementType: types.StringType,
			},
			"prefixes": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"prefix_length_v4": schema.Int64Attribute{
				Description: "prefix length",
				Optional:    true,
				Computed:    true,
			},
			"public_ip": schema.StringAttribute{
				Description: "public IP address",
				Computed:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "The project UUID.",
				Required:    true,
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
	type NetworkV0 struct {
		ID             types.String   `tfsdk:"id"`
		Name           types.String   `tfsdk:"name"`
		NameServers    types.List     `tfsdk:"nameservers"`
		Prefixes       types.List     `tfsdk:"prefixes"`
		PrefixLengthV4 types.Int64    `tfsdk:"prefix_length_v4"`
		PublicIp       types.String   `tfsdk:"public_ip"`
		ProjectID      types.String   `tfsdk:"project_id"`
		Timeouts       timeouts.Value `tfsdk:"timeouts"`
	}

	var oldState NetworkV0

	diags := req.State.Get(ctx, &oldState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ns, ds := types.SetValueFrom(ctx, types.StringType, oldState.NameServers.Elements())
	if ds.HasError() {
		resp.Diagnostics.Append(ds...)
		return
	}

	newState := Network{
		ID:             oldState.ID,
		Name:           oldState.Name,
		NameServers:    ns,
		Prefixes:       types.List{},
		PrefixLengthV4: types.Int64{},
		PublicIp:       types.String{},
		ProjectID:      types.String{},
		Timeouts:       timeouts.Value{},
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}
