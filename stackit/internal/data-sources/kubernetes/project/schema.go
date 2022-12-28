package project

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Schema returns the terraform schema structure
func (r *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource enables STACKIT Kubernetes Engine (SKE) in a project",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "kubernetes project ID",
				Computed:    true,
			},

			"project_id": schema.StringAttribute{
				Description: "the project ID that SKE will be enabled in",
				Required:    true,
			},
		},
	}
}
