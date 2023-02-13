package instance

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
	ID         types.String `tfsdk:"id"`
	InstanceID types.String `tfsdk:"instance_id"`
	ProjectID  types.String `tfsdk:"project_id"`
	Username   types.String `tfsdk:"username"`
	Database   types.String `tfsdk:"database"`
	Host       types.String `tfsdk:"host"`
	Port       types.Int64  `tfsdk:"port"`
	Role       types.String `tfsdk:"role"`
}

// Schema returns the terraform schema structure
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Data source for MongoDB Flex user\n%s",
			common.EnvironmentInfo(d.urls),
		),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the resource ID",
				Required:    true,
			},
			"instance_id": schema.StringAttribute{
				Description: "the mongo db flex instance id.",
				Required:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID the instance runs in. ",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
			},
			"username": schema.StringAttribute{
				Description: "Specifies the user's username",
				Computed:    true,
			},
			"database": schema.StringAttribute{
				Description: "Specifies the database the user can access",
				Computed:    true,
			},
			"host": schema.StringAttribute{
				Description: "Specifies the allowed user hostname",
				Computed:    true,
			},
			"port": schema.Int64Attribute{
				Description: "Specifies the port",
				Computed:    true,
			},
			"role": schema.StringAttribute{
				Description: "Specifies the role assigned to the user, either `readWrite` or `read`",
				Computed:    true,
			},
		},
	}
}
