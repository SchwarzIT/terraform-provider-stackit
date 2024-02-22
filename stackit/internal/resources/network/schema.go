package network

import (
	"context"
	"fmt"
	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Network is the schema model
type Network struct {
	Name           types.String   `tfsdk:"name"`
	NameServers    []types.String `tfsdk:"nameservers"`
	NetworkID      types.String   `tfsdk:"networkId"`
	Prefixes       []types.String `tfsdk:"prefixes"`
	PrefixLengthV4 types.Int64    `tfsdk:"prefixLengthV4"`
	PublicIp       types.String   `tssdk:"publicIp"`
	ProjectID      types.String   `tfsdk:"project_id"`
}

// Schema returns terraform schema structure
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Manages STACKIT network\n%s",
			common.EnvironmentInfo(r.urls),
		),
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "the name of the network",
				Required:    true,
				Validators: []validator.String{
					validate.NetworkName(),
				},
			},
			"nameservers": schema.ListAttribute{
				Description: "List of DNS Servers/Nameservers.",
				Required:    false,
				ElementType: types.StringType,
				Validators: []validator.List{
					validate.NameServers(),
				},
			},
			"networkId": schema.StringAttribute{
				Description: "The ID of the network",
				Required:    true,
			},
			"prefixes": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					validate.Prefixes(),
				},
			},
			"prefixLengthV4": schema.Int64Attribute{
				Description: "prefix length",
				Required:    true,
				Validators: []validator.Int64{
					validate.PrefixLengthV4(),
				},
			},
			"publicIp": schema.StringAttribute{
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
		},
	}
}
