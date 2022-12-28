package credentialsgroup

import (
	"context"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit/pkg/validate"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Schema returns the terraform schema structure
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for Object Storage credential groups",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "the credential group ID",
				Required:    true,
			},

			"object_storage_project_id": schema.StringAttribute{
				Description: "The ID returned from `stackit_object_storage_project`",
				Required:    true,
				Validators: []validator.String{
					validate.ProjectID(),
				},
			},

			"name": schema.StringAttribute{
				Description: "the credential group's display name",
				Computed:    true,
			},

			"urn": schema.StringAttribute{
				Description: "credential group URN",
				Computed:    true,
			},
		},
	}
}
