package user

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

// User is the schema model
type User struct {
	ID          types.String `tfsdk:"id"`
	ProjectID   types.String `tfsdk:"project_id"`
	InstanceID  types.String `tfsdk:"instance_id"`
	Description types.String `tfsdk:"description"`
	Write       types.Bool   `tfsdk:"writable"`
	Username    types.String `tfsdk:"username"`
}

// Schema returns the terraform schema structure
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Manages Secrets Manager instances\n%s",
			common.EnvironmentInfo(d.urls),
		),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the user UUID.",
				Computed:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "The project UUID.",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
			},
			"instance_id": schema.StringAttribute{
				Description: "Specifies the instance UUID.",
				Required:    true,
				Validators: []validator.String{
					validate.UUID(),
				},
			},
			"username": schema.StringAttribute{
				Description: "Specifies the user name.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Specifies the description of the user. Changing this value requires the resource to be recreated.",
				Computed:    true,
			},
			"writable": schema.BoolAttribute{
				Description: "Specifies if the user can write secrets. `false` by default.",
				Optional:    true,
			},
		},
	}
}
