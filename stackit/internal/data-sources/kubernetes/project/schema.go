package project

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Schema returns the terraform schema structure
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for STACKIT Kubernetes Engine (SKE) project",
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
