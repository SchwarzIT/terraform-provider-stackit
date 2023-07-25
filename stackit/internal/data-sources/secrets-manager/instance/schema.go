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

// Instance is the schema model
type Instance struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	ProjectID types.String `tfsdk:"project_id"`
	Frontend  types.String `tfsdk:"frontend_url"`
	API       types.String `tfsdk:"api_url"`
	ACL       types.Set    `tfsdk:"acl"`
}

// Schema returns the terraform schema structure
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Manages Secrets Manager instances\n%s",
			common.EnvironmentInfo(d.urls),
		),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Specifies the resource UUID.",
				Required:    true,
				Validators: []validator.String{
					validate.UUID(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "The project UUID.",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Specifies the instance name.",
				Computed:    true,
			},
			"frontend_url": schema.StringAttribute{
				Description: "Specifies the frontend for managing secrets.",
				Computed:    true,
			},
			"api_url": schema.StringAttribute{
				Description: "Specifies the API URL for managing secrets.",
				Computed:    true,
			},
			"acl": schema.SetAttribute{
				Description: "Specifies the ACLs for the instance.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}
