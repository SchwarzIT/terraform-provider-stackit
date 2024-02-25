package network

import (
	"context"
	"fmt"
	"github.com/SchwarzIT/community-stackit-go-client/pkg/services/iaas-api/v1alpha"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Network is the schema model
type Network struct {
	ID             types.String   `tfsdk:"id"`
	Name           types.String   `tfsdk:"name"`
	NameServers    types.List     `tfsdk:"nameservers"`
	NetworkID      types.String   `tfsdk:"network_id"`
	Prefixes       types.List     `tfsdk:"prefixes"`
	PrefixLengthV4 types.Int64    `tfsdk:"prefix_length_v4"`
	PublicIp       types.String   `tfsdk:"public_ip"`
	ProjectID      types.String   `tfsdk:"project_id"`
	Timeouts       timeouts.Value `tfsdk:"timeouts"`
}

// Schema returns terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Manages STACKIT network\n%s",
			common.EnvironmentInfo(r.urls),
		),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the resource ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "the name of the network",
				Required:    true,
				Validators: []validator.String{
					validate.StringWith(iaas.ValidateNetworkName, "validate network name"),
				},
			},
			"nameservers": schema.ListAttribute{
				Description: "List of DNS Servers/Nameservers.",
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					validate.NameServers(),
				},
			},
			"network_id": schema.StringAttribute{
				Description: "The ID of the network",
				Computed:    true,
			},
			"prefixes": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					validate.Prefixes(),
				},
			},
			"prefix_length_v4": schema.Int64Attribute{
				Description: "prefix length",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(25),
				Validators: []validator.Int64{
					validate.PrefixLengthV4(),
				},
			},
			"public_ip": schema.StringAttribute{
				Description: "public IP address",
				Computed:    true,
				Validators: []validator.String{
					validate.PublicIP(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "The project UUID.",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
			},
			"timeouts": common.Timeouts(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}
}
