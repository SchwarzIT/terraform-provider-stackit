package network

import (
	"context"
	"fmt"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Network is the schema model
type Network struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	NameServers    types.List   `tfsdk:"nameservers"`
	NetworkID      types.String `tfsdk:"network_id"`
	Prefixes       types.List   `tfsdk:"prefixes"`
	PrefixLengthV4 types.Int64  `tfsdk:"prefix_length_v4"`
	PublicIp       types.String `tfsdk:"public_ip"`
	ProjectID      types.String `tfsdk:"project_id"`
}

// Schema returns terraform schema structure
func (r *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Data source for STACKIT network\n%s",
			common.EnvironmentInfo(r.urls),
		),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the Network ID.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "the name of the network",
				Required:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "The project UUID.",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
			},

			// read only
			"nameservers": schema.ListAttribute{
				Description: "List of DNS Servers/Nameservers.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"network_id": schema.StringAttribute{
				Description: "The ID of the network",
				Computed:    true,
			},
			"prefixes": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"prefix_length_v4": schema.Int64Attribute{
				Description: "prefix length",
				Computed:    true,
			},
			"public_ip": schema.StringAttribute{
				Description: "public IP address",
				Computed:    true,
			},
		},
	}
}
